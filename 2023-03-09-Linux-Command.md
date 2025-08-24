---
layout: post
title: 学习笔记｜一些常见的Linux命令
categories: [Study Notes]
description: 日常比较经常使用到的一些Linux命令
keywords: Linux，centOS
mermaid: false
sequence: false
flow: false
mathjax: true
mindmap: false
mindmap2: false
---

​	记录了日常使用中，配置网络、文件目录、账号管理、系统管理、文件管理、网络管理、磁盘管理等命令。

## 配置网络

**（1）编辑ens33网卡**

```cmd
vi etc/sysconfig/network-scripts/ifcfg-ens33
```

**（2）网卡内容**

```cmd
TYPE=Ethernet

PROXY_METHOD=none

BROWSER_ONLY=no

BOOTROTO=static

IPADDR=192.168.0.0 #ip

NETMASK=255.255.255.0 #子网掩码

GATEWAY=192.168.0.0 #网关

DEFROUTE=yes

IPV4_FAILURE_FATAL=no

IPV6INIT=yes

IPV6_AUTOCONF=yes

IPV6_DEFROUTE=yes

IPV6_FAILURE_FATAL=no

IPV6_ADDR_GEN_MOD=stable-privacy

NAME=ens33

UUID=

DEVICE=ens33

ONBOOT=true
```

# 文件和目录

​	根目录/

​	系统中的配置文件 etc

​	系统预设执行文件的放置目录 usr/bin usr/sbin usr

​	程序运行日志的存放目录 var

​	切换目录 cd

​	展示目录文件 ls-l

 

**时间同步：**设置 选项 VMware tools 开启时间同步 date 查看当前系统时间

 

**常用特殊符号**

​	*代替单个或多个字符串

​	？代替一个字符串

​	[a-b]用于指定范围

​	!排除符号，用来指定被屏蔽显示内容的部分

​	`命令提大夫

​	\#注释符号

 

**重定向**

​	\>将输入的信息直接写入目标文件或设备中，并覆盖掉之前的内容

​	\>>以追加的方式写入

​	<输入重定向符号

 

**管道 |**

将前一个命令的输出作为输入，一直到最后一个命令。

# 账号管理

**（1）用户管理命令**

​	权限必须要在管理员权限下，切换用户 su

​	创建用户 useradd (选项) 用户名

​	用户口令 passwd (选项) 用户名  设置密码时，不能是一个回文。

​	修改用户 usermod 选项 用户名

​	删除用户 userdel (选项) 用户名

**（2）用户组管理命令**

​	创建用户组 groupadd (选项) 用户名

​	修改用户组 groupmod (选项) 用户名

​      	-g 将组id改为GID  -h 显示帮助信息  -n改名为NEW_GROUP

​    	  -o 允许使用重复的GID -p 将密码改为PASSWORD -R chroot到的目录

​	查询用户所属组 groups 用户名

​	删除用户组： groupdel 用户组名

**（3）管理用户组内成员**

​	gpasswd -a user group

​	-a 添加用户到组

​	-d 删除

​	-A 指定管理员

​	-M 指定组成员

​	-r 删除密码

​	-R 限制用户登陆组

​	展示组成员 grep /etc/group

# 系统管理命令

**（1）系统管理相关命令**

​	① date [参数选项] 

​	-d<字符串> 显示字符串所指的日期与时间，字符串前后必须加上双引号

​	-s<字符串> 根据字符串来设置日期与时间，字符串前后必须加上双引号

​	-u 显示GMT 英国格林威治时间

​	--help 在线帮助

​	--version 显示版本信息

​	② logname [参数选项] 显示登陆账号的信息

​	--help 帮助

​	--version 显示版本信息

​	③ 切换用户 su

​	-c 切换用户执行命令，执行完毕之后再变回原来的使用者

​	④ id [-g][--help][--version][用户名称] 查看当前用户的详细信息

​	⑤ sudo [参数选项] 提高普通用户的操作权限

​	-I 显示自己的权限

​	-u sudo -u root ls 使用root用户执行指令

​	-v 超出N分钟没有使用，要求再次输入密码

**（2）系统进程相关命令**

​	①top命令 实时显示process的动态（任务管理器）

​	-d 改变显示的更新速度

​	-q 没有任何延迟的显示速度

​	-c 切换显示模式，只显示执行档的名称/显示完整路径和名称

​	-S 累计模式 累计子行程的CPU time

​	s 安全模式

​	-p PID 显示指定进程的信息

​	结束监控的快捷键q

​	②ps命令 显示当前process进程信息

​	-A i按时系统中所有进程信息

​	-ef 显示系统中所有进程完整信息

​	-u 用户名 显示指定用户的进程信息

​	③kill命令 终端执行中的程序

​	-l<信息编号>： 若不加信息编号选项，-l参数会列出全部的信息名称

​	直接加程序的PID或者PGID，也可以是工作编号

​	-u 用户名 杀死这个用户中所有进程 killall -u 用户名

​	-9 强制杀死进程

​	kill -9 $(ps -ef | grep 用户名) 杀死用户所有进程

**（3）关机重启命令**

​	①shutdown 关机

​	-t seconds 设定在几秒后进行关机程序

​	-k 不会真的关机，警告讯息传送给所有使用者

​	-r 关机后重新开机

​	-h 关机后停机 shutdown -h now

​	-n 强行杀掉所有程序关机

​	-c 取消目前进行中的关机动作

​	-f 关机时不进行fcsk动作（检查Linux系统）

​	-F 关机时强迫进行fsck动作

​	time  设定关机的事件

​	-message 传送给所有使用者的警告讯息

​	②reboot 重启命令

**（4）防火墙设置**

​	iptables服务

​	sudo sysv-rc-conf iptables on

​	sudo sysv-rc-conf –list iptables

​	sudo sysv-rc-conf –level 2345 iptables

**（5）系统管理的其他命令**

​	①Who命令 显示当前登陆系统的用户

​	-H 显示明细标题信息

​	②timedatectl命令 校正服务器时间、时区

​	status 显示系统的当前时间和日期

​	list-timezones 查看所有可用的时区

​	set-timezone “Asia/Shanghai” 设置本地时区

​	set-ntp false 禁用时间同步

​	set-time “2019-03-11 20:45:00” 设置时间

​	③Clear 清除屏幕

# 文件管理和编辑命令

**（1）目录管理命令**

​	①ls 列出目录

​	-a 显示所有文件或目录（包含隐藏）

​	-d 仅列出目录本身，不列出目录内的文件数据

​	-l 长数据串列出，包含文件的属性与权限等数据

​	权限 d：目录或文件夹 -表示文件  

​	属主 属组 大小 最后一次访问时间 名称

​	②pwd 查看当前所在的目录

​	-P 效果一样

​	③cd [相对路径或绝对路径] 切换目录

​	.. 可以回退

​	④mkdir命令  创建目录/文件夹

​	mkdir [-p] 文件夹名字 创建多级文件夹

​	⑤rmdir 命令  删除文件夹

​	只能删除空文件夹

​	rkdir [-p] 文件夹名字  删除多级文件夹后，若上一级为空则全删除

​	⑥rm [选项] 文件/目录  删除目录或者文件

​	-I 删除前逐一询问

​	-f 只读文件直接删除，无需逐一确认

​	-r 将目录及以下文件均删除

​	⑦touch 创建文件

​	⑧cp命令  文件复制

​	cp [选项] 数据源 目的地

​	-a 复制目录时使用，保留链接、文件属性，及目录下所有内容

​	-d 复制时百六链接

​	-f 覆盖已经存在的目标文件而不给出提示

​	-I 覆盖文件前提示

​	-p 同时复制修改时间和访问权限

​	-r/R 复制目录下所有子目录和文件

​	-l 不复制文件，只生成链接文件

​	\* 表示所有文件

​	⑨mv命令 移动文件夹或文件 改名

​	mv [选项] 数据源 目的地

​	-i 同名文件询问是否覆盖

​	-f 直接覆盖不提示

**（2）文件基本属性**

​	可读r 可写w 可执行x 没有当前权限-

​	第1位：d 目录  -文件   l为链接文档

​	第2-4位：属主权限

​	第5-7位：属组权限

​	第8-10位： 其他用户权限

​	①chgrp命令 更改属组

​	-v 增加提示语句

​	chgrp -v root aaa 将aaa的属组改为root

​	②chown命令 更改属组和属主

​	-R 处理指定目录及目录下所有文件

​	chown -r root:root aaa 讲aaa文件夹和里面所有的属主和属组改为root

​	③chmod命令 修改属主、属组、其他用户权限

​	数字权限 r4 w2 x1 无权限0

​	rwx=4+2+1=7

​	chmod [参数选项] 数字权限 文件或目录

​	-R 对所有进行变更权限

​	符号权限

​	chmod [参数选项] 符号权限 文件或目录

​	user 属主权限 u

​	group 属组权限 g

​	others 其他权限 o

​	全部身份 a

​	+加入 -出去 =设定

​	chmod u=rwx，g=rx，o=r a.txt

**（3）文件管理**

​	①touch命令

​	touch [参数选项] 文件名 不存在就创建，存在修改时间属性

​	查看文件属性 stat 文件名

​	touch a{1..10}.txt 批量创建空文件

​	②vi/vim编辑器

​	vi 只能编辑文本内容，不能对字体段落进行排版 不支持鼠标操作 没有菜单 只有命令

​	vim 从vi发展出来的文本编辑器 优于代码补全、编译及错误跳转等编程功能

​	vi/vim的三种模式：阅读（命令模式），编辑（编辑模式），保存（末行模式）

​	i：从命令模式进入编辑模式

​	Esc：从编辑模式退出命令模式

​	“：”冒号：从命令模式进入末行模式

​	两下Esc：从末行模式进入命令模式

​	1）  打开和新建文件

​	vim 文件名 存在则打开 不存在则创建临时文件，保存后会新建

![image-20240903104405752](/images/posts/2023-03-09-Linux-Command.assets/image-20240903104405752.png)

​	1）  末行模式保存文件

![image-20240903104414629](/images/posts/2023-03-09-Linux-Command.assets\image-20240903104414629.png)

​	Ctrl+u 向上翻半页。

​	Ctrl+d 向下翻页。

​	Ctrl+g   显示光标所在位置的行号和文件的总行数。

​	nG 光标跳到文件的第n行行首。

​	G  光标跳到文件最后一行。

​	:5回车  光标跳到第5行。

​	:n回车  光标跳到第n行。

​	0  光标跳到当前行的行首。

​	$  光标跳到当前行的行尾。

​	w  光标跳到下个单词的开头。

​	b  光标跳到上个单词的开头。

​	e  光标跳到本单词的尾部。

​	x   每按一次，删除光标所在位置的一个字符。

​	nx  如"3x"表示删除光标所在位置开始的3个字符。

​	dw 删除光标所在位置到本单词结尾的字符。

​	D  删除本行光标所在位置后面全部的内容。

​	dd  删除光标所在位置的一行。

​	ndd 如"3dd"表示删除光标所在位置开始的3行。

​	yy  将光标所在位置的一行复制到缓冲区。

​	nyy 将光标所在位置的n行复制到缓冲区。

​	p  将缓冲区里的内容粘贴到光标所在位置。

​	r   替换光标所在位置的一个字符 replace。

​	R  从光标所在位置开始替换，直到按下"Esc"。

​	cw 从光标所在位置开始替换单词，直到按下"Esc"。

​	u  撤销命令，可多次撤销。

​	J  把当前行的下一行接到当前行的尾部。

​	/abcd 在当前打开的文件中查找“abcd”文本内容。

​	n   查找下一个。

​	N   查找上一下。

​	.  重复执行上一次执行的vi命令。

​	~  对光标当前所在的位置的字符进行大小写转换。

​	列操作

​	Ctrl+V  光标上或下 大写的I 输入内容 Esc

​	③其他命令

![image-20240903104442168](/images/posts/2023-03-09-Linux-Command.assets\image-20240903104442168.png)

​	cat -N 显示行号

![image-20240903104447496](/images/posts/2023-03-09-Linux-Command.assets\image-20240903104447496.png)

​	④tail [参数选项] 文件 查看文件的最后部分

![image-20240903104452245](/images/posts/2023-03-09-Linux-Command.assets\image-20240903104452245.png)

​	tail -f 动态显示 ctrl+c退出

​	⑤grep [参数选项] 关键字 文件  根据关键词，搜索文本文件内容

![image-20240903104457881](/images/posts/2023-03-09-Linux-Command.assets\image-20240903104457881.png)

​	-c 查找进程个数

​	⑥vim定位行

​	vim 文件名 行数 查看文件并定位到具体行数

​	异常处理 异常退出时，在磁盘上可能保存有交换文件 删除后可正常访问

​	⑦echo 命令 

​	echo字符串 展示文本

​	echo 字符串 > 文件名 将字符串写到文件中（覆盖文件中内容）

​	echo 字符串 >> 文件名 将字符串写到文件中（不覆盖文件中内容）

​	cat 不存在的目录 &>> error.log 将命令的失败结果追加到error.log文件的后面

​	⑧软连接（快捷方式）

​	ln -s 文件路径 快捷方式名

​	⑨find 查找命令

​	find [参数选项] <指定目录> <指定条件> <指定内容> 在指定目录下查找文件

![image-20240903104524294](/images/posts/2023-03-09-Linux-Command.assets\image-20240903104524294.png)

​	示例：

​	从/tmp目录开始搜索，把全部的*.c文件显示出来。

​	find /tmp -name *.c -print

​	从当前工作目录开始搜索，把全部的*.c文件显示出来。

​	find . -name *.c -print

**（4）压缩命令 gzip gunzip**

​	① gzip [参数选项] [文件] 压缩文件 *压缩所有文件

​	-dv 解压文件并列出详细信息

​	② gunzip [参数选项] [文件] 解压文件

​	③ tar [必要参数] [选择参数] [文件] 打包、压缩和解压（文件/文件夹）

![image-20240903104534451](/images/posts/2023-03-09-Linux-Command.assets\image-20240903104534451.png)

​	④ zip命令

​	zip [必要参数][选择参数][文件] 压缩

​	-q 不显示指令执行过程

​	-r 递归处理，将指定目录下的所有文件和子目录一并处理

​	Zip -q -r aaa.zip aaa

​	⑤ unzip命令

​	⑥rpm命令 管理程序

​	只能安装已经下载到本地机器上的rpm包

**（5） yum命令** 

​	查找安装下载卸载命令

![image-20240903104545233](/images/posts/2023-03-09-Linux-Command.assets\image-20240903104545233.png)

​	-y 若需确认 一致确定

​	更改yum源

![image-20240903104551743](/images/posts/2023-03-09-Linux-Command.assets\image-20240903104551743.png)

**（6）统计文本文件的行数、单词数和字节数**

​	wc 文件名

​	示例：

​	统计当前工作目录处book2*.c文件的行数、单词数和字节数。

​	wc book2*.c

# 网络管理

**（1）Ifconfig命令**  

​	用于显示或配置网络设备的命令

​	网卡 down 关闭网卡

​	网卡 up 打开网卡

​	网卡 ip 配置ip地址

​	网卡 ip +netmask 掩码 配置ip地址和子网掩码

​	域名解析测试nslookup

​	测试网络状态 ping

​	traceroute命令查看数据包经过的路由路径

**（2）ping命令 检测是否与主机联通**

​	ping 网址 检测联通

​	ping -c 2 网址 联通次数

**（3）netstat命令 显示网络状态**

​	显示活动的tcp连接、计算机侦听端口、以太网统计信息等

​	-a 显示所有连线中的socket

​	-i 显示网卡列表

**（4）网络配置工具ip**

​	Ip [OPTIONS] OBJECT [COMMAND[ARGUMENTS]]

​	-V打印iproute信息

​	-r 将ip地址转换成域名

​	-s：输出详细结果

​	Command指令：add添加delete删除lis/show列表help帮助

**（5）ftp访问命令**

​	ftp [-dignv] [主机名称ip地址]

​	-d显示详细指令执行过程

​	-g关闭本地主机文件名称

​	-I 关闭互动模式，不询问任何问题

​	-n不使用自动登录

​	-v显示指令执行过程。

​	ftp内部指令

​	ls列出当前目录

​	cd在远程机上改变工作目录

​	lcd在本地机上改变工作目录

​	close终止当前的ftp会话

​	hash每次传输玩数据缓冲区中的数据后就显示一个#号

​	get(mget)从远程机传送指定文件到本地机

​	put(mput)从本地机传送指定文件到远程机

​	quit断开与远程机的连接，并退出ftp

**（6）wget命令行下载工具**

​	wget [参数] [URL地址]

​	-O 把文档写到FILE文件中

​	-c 接着下载没下载完的文件

# 磁盘管理方法

**（1）磁盘分区命令 fdisk[参数]**

​	-b<分区大小>：指定每个分区的大小

​	-l：列出指定的外围设备的分区情况。

​	-s<分区编号>：将指定的分区大小输出到标准输出上，单位为区块。

​	-u：搭配”-l”参数列表，会用分区数目取代柱面数目，来表示每个分区的起始地址。

**（2）建立文件系统命令 mkfs**

​	mkfs [参数] [-t 文件系统类型] [分区名称]

​	fs：指定建立文件系统时的参数。

​	-v：显示版本信息和详细的使用方法。

​	-V：显示版本信息和简要的使用方法，主要用于测试。

**（3）设置交换分区命令 mkswap**

​	Mkswap [参数] [设备或分区名称] [交换分区大小]

​	-c：在建立交换分区前，检查块设备是否有损坏，如果有，打印出来。

​	-f：在SPARC机器上建立交换分区时使用

​	-v0：建立旧式的交换分区。

​	-v1：建立新式的交换分区

**（4）显示磁盘信息命令 df**

​	df [参数] [设备或区块]

​	-a：包含全部的文件系统

**（5）显示目录的容量 du**

​	du [参数] [目标目录]

​	-a：显示目录中所有文件的大小

**（6）挂载与卸载分区**

​	mount命令操作

​	mount [参数] -t [类型] [设备名称] [目的地目录]

​	-V：显示程序版本

​	-h：显示辅助信息

​	-a：将/etc/fstab中定义的所有档案系统挂上

​	-F：为每一个mount动作产生一个行程负责执行，加快挂载速度。

​	-f：用来除错，不执行实际动作，而是模拟。

# Git命令

**（1）初始化**

​	使用git init对当前目录进行初始化，会生成./git隐藏文件。

**（2）增加与提交**

​	git add “文件” 将目标文件添加到git仓库

​	git commit -m “本次提交说明” 提交本次更改

**（3）版本管理**

​	git status命令可以让我们时刻掌握仓库当前的状态

​	git diff顾名思义就是查看difference

​	git log命令显示从最近到最远的提交日志 加上--pretty=oneline参数简化显示信息

​	git reset --hard HEAD~n 回退n版

​	git reset --hard 1094a 也可以直接输入序号，回退指定版本

​	git reflog用来记录你的每一次命令

​	git checkout -- file 让这个文件回到最近一次git commit或git add时的状态

​	git rm 在版本库中删除对应文件

**（4）远程仓库**

​	git remote add origin [git@github.com:michaelliao/learngit.git](mailto:git@github.com:michaelliao/learngit.git)

​	git push origin master 推送到仓库

​	git clone [git@github.com:michaelliao/gitskills.git](mailto:git@github.com:michaelliao/gitskills.git) 克隆库到本地

**（5）其他**

​	tmux new -s roclinux -s是 session 的缩写，顾名思义，我们启动了一个全新的 tmux 会话（tmux session），并且把这个会话起名叫作 roclinux

​	按 Ctrl-B 组合键+n 切换窗口

​	按 Ctrl-B 组合键+c 增加窗口

​	按 Ctrl-B 组合键+d 退出窗口

​	Tmux ls 展示存在的窗口

​	Tmux a -t reclinux 返回窗口

​	启动一个新的tmux会话：

​	tmux

​	以指定名称启动一个新的tmux会话：

​	tmux new -s session_name

​	列出现有的tmux会话：

​	tmux ls

​	进入指定的tmux会话：

​	tmux attach-session -t session_name

​	以只读模式进入指定的tmux会话：

​	tmux attach-session -t session_name -r

​	分离当前的tmux会话：

​	tmux detach

​	在当前tmux会话中创建一个新的窗口：

​	Ctrl-b c

​	关闭当前的tmux窗口：

​	Ctrl-d

​	将当前的tmux窗口重命名：

​	Ctrl-b ,

​	在不同的窗口之间切换：

​	Ctrl-b n

​	（下一个窗口）、

​	Ctrl-b p

​	（上一个窗口）

​	在tmux窗口中水平分割窗格：

​	Ctrl-b %

​	在tmux窗口中垂直分割窗格：

​	Ctrl-b "

​	在tmux窗格之间移动：

​	Ctrl-b ↑

​	Ctrl-b ↓

​	Ctrl-b ←

​	Ctrl-b →

​	改变窗格的大小：

​	Ctrl-b Alt-↑

​	Ctrl-b Alt-↓

​	Ctrl-b Alt-←

​	Ctrl-b Alt-→

​	显示当前tmux会话状态栏：

​	Ctrl-b t

​	显示所有tmux快捷键：

​	Ctrl-b ?
