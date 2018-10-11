package main
//import include
import(
	"fmt"
	"io"
	"os"
	"os/exec"
	"bufio"
//	"syscall"
	flag "github.com/spf13/pflag"
)
//types struct
type sp_args struct
{
	start_page int
	end_page int
	in_filename string
	page_len int
	page_type string

	print_dest string
}
//globals
var sa = new(sp_args) 
//program name,for error messages
var progname string

func main() {
	progname = os.Args[0]

	flag.IntVarP(&sa.start_page, "s", "s", -1, "start page number(-1)")
	flag.IntVarP(&sa.end_page, "e", "e", -1, "end page number(-1)")
	flag.IntVarP(&sa.page_len, "l", "l", 72, "lines per page(72)")
	flag.StringVarP(&sa.page_type, "f", "f", "l", "form->l(default)/f")
	flag.StringVarP(&sa.print_dest, "d", "d", "", "destination")
	// fei bi xu xuan xiang zhi
	flag.Lookup("f").NoOptDefVal = "f"


	flag.Usage = usage

	flag.Parse()

	process_args()
	process_input()
}



//process_args() 
func process_args() {
	
	/* check the command-line arguments for validity */
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "\n%s: not enough command-line arguments\n", progname)
		flag.Usage()
		os.Exit(1)
	} 
	if os.Args[1] != "-s" || os.Args[3] != "-e" {
		fmt.Fprintf(os.Stderr, "\n%s: 1st arg should be -s start_page and 2nd arg should be -e\n", progname)		
		flag.Usage()
		os.Exit(2)
	}
	if sa.start_page < 1 {
		fmt.Fprintf(os.Stderr, "invalid start page %v\n", sa.start_page)
		flag.Usage()
		os.Exit(3)
	}
	if sa.end_page < 1 || sa.end_page < sa.start_page {
		fmt.Fprintf(os.Stderr, "invalid end page %v\n", sa.end_page)
		flag.Usage()
		os.Exit(4)
	} 
		if flag.NArg() > 0 {
			
			_, err := os.Stat(flag.Args()[0])
			if  err != nil && os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "\n%s: input file \"%s\" does not exist\n",progname, flag.Args()[0]);
				os.Exit(5);
			}
			sa.in_filename = flag.Args()[0]

		}
	

}
//process_input()

func process_input(){
	var fin *os.File
	var fout io.WriteCloser
	bufFin := bufio.NewReader(fin)
	cmd := &exec.Cmd{}
	var page_ctr int
	var line_ctr int
	var err error

  // input


	if sa.in_filename != "" {
		fin, err = os.Open(sa.in_filename)
		
		if err != nil {
			fmt.Fprintf(os.Stderr, "\n%s: could not open input file \"%s\"\n",progname, sa.in_filename)
			os.Exit(6)
		}
		defer fin.Close()
	} else {
		fin = os.Stdin
	}

  //output
  if sa.print_dest == ""{
  	fout = os.Stdout
  }else {
  	cmd = exec.Command("lp", "-d", sa.print_dest)
  	fout, err = cmd.StdinPipe()
		if err != nil {
			fmt.Fprintf(os.Stderr, "\n%s: can't open pipe to \"lp -d%s\"\n",progname, sa.print_dest)
			os.Exit(7)
		}
		cmd.Start()
  }

  // count by line
  if sa.page_type == "l" {
  	line_ctr = 0
  	page_ctr = 1
  	for{
  		line, rerror := bufFin.ReadString('\n')
  		if rerror != nil {break}

  		line_ctr++

  		if line_ctr > sa.page_len {
			line_ctr=1
  			page_ctr++
		}
  		if page_ctr > sa.end_page {break}

  		if page_ctr >= sa.start_page && page_ctr <= sa.end_page {
				_, werror := fout.Write([]byte(line))
  			if werror != nil{ os.Exit(4) }
  		}


  	}

  } else{//count by /n
  	page_ctr = 1
  	for{

  		line, rerror := bufFin.ReadString('\n')
  		if rerror != nil {break}

  		if(page_ctr >= sa.start_page)&&(page_ctr <= sa.end_page) {
  			
  			_, werror := fout.Write([]byte(line))
  			if werror != nil{ os.Exit(4) }

  		}
  		page_ctr++
  	}
  }

  cmd.Wait()
  defer fout.Close()

  if page_ctr < sa.start_page {
  	fmt.Fprintf(os.Stderr, "start_page (%d) > total pages (%d), no output written\n", sa.start_page, page_ctr)
  } else if page_ctr < sa.end_page {
  	fmt.Fprintf(os.Stderr, "end_page (%d) > total pages (%d), less output than expected\n", sa.end_page, page_ctr)
  }
}

//usage() 

func usage() {
	fmt.Fprintf(os.Stderr,
			"\nUSAGE: %s -s start_page -e end_page [ -f | -l lines_per_page ]" + 
			" [ -d dest ] [ in_filename ]\n", progname)
		flag.PrintDefaults()
}

// EOF
