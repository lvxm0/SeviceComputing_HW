<h1>设计思路</h1>
1. 按照selpg.c的代码逻辑：
  导入所需要的包，定义结构体，并定义实现下列函数：process_args()	process_input() usage()<br>
<h1>使用说明</h1>
先在终端运行 $ go get github.com/spf13/pflag<br>
然后运行 $ go build selpg.go<br> 
运行selpg $ ./selpg<br>

<h1>测试结果</h1>
1. 测试数据输入为100行，每一行为对应行号。先测试不加参数的命令行。<br>

2.测试命令 $ ./selpg -s 1 -e 2 input.txt 
![](https://github.com/lvxm0/SeviceComputing_HW/raw/master/HW1/01.PNG)
