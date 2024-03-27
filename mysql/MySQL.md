# 数据库

********

## 数据库的简单介绍：

### 数据库概述：

> 数据库：
>
> ​	英文单词DataBase，简称DB。按照一定格式存储数据的一些文件的组合或者集合
>
> 数据可管理系统：
>
> ​	DataBaseManagment, 简称DBMS。专门用来管理数据库的，对数据库中的数据进行增删改查。
>
> ​	常见的数据库系统：MySQL，Oracle，MS SqlServer，DB2，sybase……

> SQL：结构化查询语言
>
> ​	程序员通过编写SQL语句，然后DSMS负责执行SQL语句，完成数据库的增删改查操作。

> 三者之间的关系：
>
> DBMS --> 执行 --> SQL --> 操作 -->DB

### 常用的数据库命令

退出MySQl

> `exit`

显示所有的数据库

> `show DataBases;`
>
> MySQL默认自带了四个数据库

选择使用数据库

> use + "库名";
>
> `use text;`

创建数据库

> create databases + "数据库名"；

查看数据库中那些表

> show tables;

查看数据库的版本号

> select version();

查看当前使用的数据库

> select database();

终止命令的输入

> /C 或者 ; 结束

查看表中的记录和字段

> desc + 表名

**注意：MySQL不见分号不执行，当输入分号时才执行语句。**

### 表(table)的理解

数据库中的最基本单元：表。

数据库当中是以表格的形式来表示数据。因为表比较直观。

任何一张表都有行和列

​	行（row）：被称为数据/ 记录

​	列（column）：被称为字段

### SQL语句的分类

> DQL:数据查询语言（凡是带有select关键字的语言）
>
> ​	select ……
>
> DML:数据操作语言（凡是对**表当中的数据进行增删改查的都是DML**）
>
> ​	insert 	增
>
> ​	delete 	删
>
> ​	update	改
>
> DDL:数据定义语言（凡是带有create，drop，alter的都是DDL）
>
> ​	**DDL主要操作的是表的结构。不是表中的数据**。
>
> ​	create 	新建
>
> ​	drop		删除
>
> ​	alter 		修改
>
> TCL:事务控制语言
>
> ​	事务提交：commit
>
> ​	事务回滚：rollback
>
> DCL:是数据控制语言
>
> ​	例如：授权：		grant。
>
> ​				撤销权限：revoke。

### 导入数据

> source + "路径"
>
> 路径中不能有中文

### 查看表结构

> 查看表中的数据
>
> ```sql
> select * from 表名;  //统一执行这个SQL语句。
> ```
>
> eg：select \* from sys;

> 查看表的基本结构
>
> ```sql
> describe缩写为desc。
> desc 表名;
> eg：desc sys;
> ```

> 查看表的详细结构
> ```sql
> show create table 表名;
> eg：show create table emp_1;
> ```

## 简单查询语句

### 查询一个字段

select 字段名 from 表名;

​	其中：select和from都是关键字

​				字段名和表名是标识符。

### 查询多个字段

使用逗号隔开

select 字段名,字段名,字段名 from 表名;

### 查询所有字段

1. 写上每个字段

   select a,b,c,d... from 表名;

2. 可以使用  \* ;

   select \* from 表名;

   缺点：\*会先转换称为所有的字段名。

   ​			效率低；

   ​			可读性差，

   ​			不建议在实际开发中使用。

### 给查询的列起别名

```sql
select 字段名,字段名 as 别名 from 表名;
as关键字可以省略
select 字段名,字段名 别名 from 表名;
```

注意：只是将显示的结果列名显示为别名，并不会改变表中的列名。

​			select语句只负责查询，永远不会进行修改操作。

起别名时，别名里面有空格的时候

`select 字段名,字段名 as 别名(dept name) from 表名;`

别名中有空格时，使用单引号或者*<u>**双引号括住**</u>*。

**注意：在所有的数据库中，字符串统一的用法是用单引号括起来，单引号是标准的用法，双引号在Oracle中用不了。但是在MySQL中可以使用。**

### 列参与数学运算

字段可以使用数学表达式参与运算

```sql
select 
	ename, sal*12 
from 
	emp;
select 字段名 字段名+数学运算 from 表名

起别名
select 
	ename, sal*12 as yearsal 
from 
	emp;
	
别名是中文时用单引号括起来。
select ename, sal*12 as '年薪' from emp;
```

### 条件查询

不查询表中的所有数据，只查询出来符合条件的数据。

语法格式：

```sql
select
	字段1，字段2，字段3……
from
	表名
where
	条件;
```

> 条件：
>
> ​	= 等于
>
> ​	`查询薪资等于800的员工姓名和编号：`
>
> ​	`select empno,ename from emp where sal = 800'`
>
> ​	<> 或 != 不等于
>
> ​	`select empno,ename from emp where sal != 800;`
>
> ​	`select empno,ename from emp where sal <> 800;`
>
> ​	< 小于
>
> ​	<=  小于等于
>
> ​	> 大于
>
> ​	 >= 大于等于
>
> ​	between ... and ... 两个值之间，等同于>= and <=
>
> ​	1,
>
> ​	`查询薪资在2450到3000之间的员工信息？包括2450和3000`
>
> ​	`select empno, ename,sal from emp where sal >= 2450 and sal <= 3000;`
>
> ​	2,
>
> ​	`select empno,ename,sal from emp where sal between 2450 and 3000;`
>
> ​	**注意：使用between……and……的时候，必须左小右大，左大	右小查不出来值。**
>
> ​	 is null 为 null  (is not null 不为空)
>
> ​	`查询那些员工的薪资或者津贴为null`
>
> ​	`select empno,ename,sal from emp where comm is null;`
>
> ​	**注意：在数据库中null不能使用等号进行衡量。需要使用is null。因为数据库中的null代表什么也没有，不是一个值，所以不饿能使用等号衡量。**
>
> ​	`查询那些员工的薪资或者津贴不为null`
>
> `select ` 
>
> ​	`	empno,ename,sal `		
>
> `from `
>
> ​	`	emp `
>
> `where `
>
> ​	`comm is not null;`
>
> ​	and 并且
>
> ​	`查询工作岗位是MANAGE并且工资大于2500的员工信息`
>
> `select `
>
> ​	`			empuo,ename,job,sal `			
>
> `	from `	
>
> ​	`emp `
>
> `where `
>
> ​	`job = 'MANAGER' and sal > 2500;`
>
> ​	or 或者
>
> `select`
>
> ​	`empno,ename,job`
>
> `from`
>
> ​	`emp`
>
> `where`
>
> ​	 `job = 'MANAGER'or'SALESMAN';`
>
> ​	**and和or同时出现的时候有没有先后顺序的问题：and的优先级比or高,想要or先执行加小括号**
>
> ​	`找出工资大于2500并且部门编号为10或者20的所有员工：`
>
> `select`
>
> ​	`*`
>
> `from`
>
> ​	`emp`
>
> `where`
>
> ​	`sal > 2500 and (deptno = 10 or deptno = 20);`
>
> ​	in 包含，相当于多个 or （not in不在这个范围中）
>
> ​	`查询工作岗位是‘MANAGER’和‘SALESMAN’`
>
> ​	`select`
>
> ​		`	empno,ename,sal`
>
> ​	`from`
>
> ​		`emp`
>
> ​	`where`
>
> ​		`job = 'MANAGER' or job = 'SALESMAN'`;
>
> ​	`select`
>
> ​		`	empno,ename,sal`
>
> ​	`from`
>
> ​		`emp`
>
> ​	`where`
>
> ​		`job in ('MANAGER','SALESMAN'); //这个不是表示800-5000，表示的是是800和5000的。`
>
> ​	not 可以取反，主要用在is 或 in 中
>
> ​		is null;
>
> ​		is not null;
>
> ​		in;
>
> ​		not in;
>
> ​	like 
>
> ​		称为模糊查询，支持%和下划线匹配
>
> ​		% 匹配任意多个字符
>
> ​		下划线：任意一个字符
>
> ​		%和__是两个特殊符号。
>
> ​	`找出名字中含有o的：`
>
> ​		`select ename from emp where ename like '%o%;'`
>
> ​	`找出名字中以t结尾的：`
>
> ​		`select ename from emp where ename like '%t;'`
>
> ​	`找出名字中第一个是k的：`
>
> ​		`select ename from emp where ename like 'k%;'`
>
> ​	`找出名字中第二个字母是A的：`
>
> ​		`select ename from emp where ename like '_A%;'`
>
> ​	`找出名字中第三个字母是R的：`
>
> ​		`select ename from emp where ename like '__R%;'`
>
> ​	`找出名字中有'_'的：//: \是转义字符。`
>
> ​		`select ename from emp where ename like '%\_%;'`

## 排序

### 查询所有员工薪资，排序

```sql
select
	ename,sal
from
	emp
order by
	sal;
//默认是升序排列。
//指定降序
select
	ename,sal
from
	emp
order by
	sal desc;
//指定升序
select
	ename,sal
from
	emp
order by
	sal asc;
```

### 两个字段排序，多个字段排序

```sql
//查询员工名字和薪资，要求按照薪资升序，如果薪资一样的话，再按照名字升序排列。
select
	ename,sal
from
	emp
order by
	sal asc, ename asc;
```

### 根据字段的位置排列

```sql
select ename,sal from emp order by 2;
//2表示第二列。是sal。不在开发中这样写，因为程序不健壮。
```

### 综合案例

找出工资在1250-3000之间的员工信息，要求按照薪资降序排列

```sql
select
	ename,sal
from
	emp;
where
	sal between 1250 and 3000
order by
	sal desc;
//关键字顺序不能改变
select
	...
from
	...
where
	...
order by
	...
程序的执行步骤：
	第一步：执行from从table中查询
	第二步：执行where筛选过滤
	第三步：执行select输出数据
	第四步：执行order by，排序总是在最后执行！
```

## 数据处理函数（单行处理函数）

### 数据处理函数又被称为单行处理函数

单行处理函数的特点是：一个输入一个输出。

和单行处理函数对应的是：多行处理函数（多个输入，一个输出）。

### 常见的单行处理函数：

> sum					求和
>
> lower				转换小写
>
> ​	`将姓名转换成小写字母显示`
>
> ​	`select lower(ename) as ename from emp;`
>
> upper				转换大写
>
> substr				取子串（substr（被截取的字符串，起始下标，截取的长度））
>
> ​	`select substr(ename, 1, 1) as ename from emp;`
>
> ​	**注意：起始下标从1开始，没有0。**
>
> ```sql
> //找出员工名字第一个字母是A的员工信息？
> 
> //第一种方式：模糊查询	
> select ename from emp where ename like 'A%';
> 
> //第二种方式：substr函数	
> select ename，from emp，where substr(ename, 1, 1) = 'A';
> 
> //首字母大写	
> select concat(upper(substr(name, 1, 1)), substr(name,2,length(name) - 1)) as result	from emp;
> ```
>
> length				取长度
>
> ​	`select length(ename) enamelength from emp;`
>
> trim					去空格
>
> ​	`//去除查询条件的空格。`
>
> ​	`select * from emp where ename = trim('  KING');`
>
> str_to_date		将字符串转换成日期
>
> date_format		格式化日期
>
> format				设置千分位
>
> case...when...then...when...then..else...end
>
> ```sql
> //当工作岗位是MANAGER的时候，工资上调10%，当工作岗位是SALESMAN的时候，工资上调50%
> //注意：不修改数据库，只是将查询结果显示为工资上调。
> select ename, job, sal as oldsal, 
> 
> (case job 
> 	when 'MANAGER' 
> then sal * 1.1 
> 	when 'SALSMAN' 
> then sal * 1.5 
> 	else sal 
> end)
> 	as newsal 
> from 
> 	emp;
> ```
>
> round				四舍五入	
>
> ​	`select 字段名 from 表名;`
>
> ​	`select 'abc'/1000 from emp;`
>
> ​	**注意：select后面可以跟某个表的字段名（可以等同于变量名），	也可以跟'字面值/字面量'（数据），**
>
> ​	`select 21000 as num from dept;`
>
> ​	`select round(1234.6756, 0) as result from emp;`
>
> ​	`//四舍五入的结果是1234。0是不保留小数，1，保留一位小数.1234.5。-1保留到十位.1230.-2，保留到百位1200`
>
> rand()				生成随机数
>
> ifnull					可以将null转换成一个具体值(如果是空的话，当作一个值来处理。)
>
> **在所有的数据库当中，只要有null参与的数学运算的结果就是NULL。**
>
> ​	`计算所有员工的年薪：`
>
> ​	`select ename, sal + comm as salcomm from emp;`
>
> ​	为了避免出现NULL值。使用ifnull函数（如果是null，当作0来处理。）
>
> ​	`select ename, (sal + ifnull(comm, 0)) * 12 as yearsal from emp;`

## 分组函数（多行处理函数）

多行处理函数的特点：输入多行，最终输出一行

5个：

1. count		计数
2. sum     求和
3. avg    平均值
4. max   最大值
5. min    最小值

**注意：分组函数在使用的时候必须先进行分组，然后才能使用。如果你没有对数据进行分组，整张表默认为一组。**

> ```sql
> 找出最高工资
> select max(sal) from emp;
> 找出最低工资
> select min(sal) from emp;
> 计算工资和
> select sum(sal) from emp;
> 计算平均工资
> select avg(sal) from emp;
> 计算员工数量总和
> select count(ename) from emp;
> ```
>
> **分组函数在使用时的注意事项：**
>
> 1，分组函数在使用的时候会自动忽略null；null不是一个值，是什么也没有。不用加null过滤。
>
> 2，分组函数中count(\*)和count(具体字段)的区别：
>
> ​	count(\*)：统计表当中的总行数。（只要有一行数据，count ++）
>
> ​						因为每一行记录中数据不可能都为null。一行数据中有一列不为null，则这行数据就是有效的。
>
> ​	count(具体字段)：表示统计该字段下所有不为null的元素的总数。
>
> 3，分组程序不能够直接使用在where子句中
>
> ```sql
> //找出比最低工资高的员工的信息。
> select ename,sal from emp where sal > min(sal);
> //表面没有问题，但运行时出现错误。
> //Invaild use of group function因为分组函数在使用的时候必须先分组，where执行的时候还没有分组。所以where后面不能出现分组函数
> ```
>
> 4，所有的分组函数可以组合在一起使用。
>
> `select sum(sal),min(sal),max(sal),count(*),avg(sal) from emp;`

## 分组查询（***\*******\*）

概述：在实际应用中，可能有这样的需求，需要先进行分组，然后对每一组的数据进行操作。这个时候需要分组查。

```sql
select	...from	...group by	...
//计算每个部门的工资和//计算每个工作岗位的平均薪资
//找出每个工作岗位的最高薪资……组合之前的关键字，分析它们的执行顺序
select	...from	...group by	...order by	...
以上关键字的顺序不能颠倒。执行顺序：	
	1，from	
	2，where	
	3，group by	
	4，select	
	5，order by select sum(sal) from emp;
这个没有分组，为什么可以执行
因为select在group by后面执行。
```

### 找出每个工作岗位的工资和

```sql
实现思路:按照工作岗位进行分组，对分组逐个求和。
select	
	job,sum(sal)
from 
	emp 
group by	
	job;
以上语句的执行顺序：先从emp表中进行查询，根据job字段进行分组，然后对每一组的数据进sum(sal);	
select	
	ename.job.sum(sal) 
from 
	emp
group by	
	job;
以上代码在oracle执行报错，在mysql中可以执行，但是毫无意义。Oracle的语法比MySQL严格
```

**结论：在一条select语句后面，如果有group by语句的话；select后面只能跟：参加分组的字段，以及分组函数。其他的一律不能跟。**

### 找出每个部门的最高薪资

```SQL
select 	
	ename,deptno,max(sal)
from 
	emp 
group by 
	deptno;	
	
select字段后面加ename没有意义，另外Oracle会报错。正确写法：
select 
	deptno,max(sal)
from 
	emp 
group by 
	deptno;
```

### 找出<u>每个部门，不同工作岗位</u>的最高薪资

```sql
两个字段联合成一个字段看
select	
	deptno,job,max(sal)
from	
	emp 
group by 
	job, deptno;	
```

### 找出每个部门的最高薪资,要求显示最高薪资大于3000的。使用having可以对分完组之后的数据进一步过滤。

> ***having不能单独使用，不能代替where having，必须和group by联合使用。***

```sql
第一步：
select 
	deptno,max(sal) 
from 
	emp 
group by 
	deptno;
第二步：
select 	
	deptno,max(sal) 
from 
	emp 
group by 
	deptno 
having	
	max(sal) > 3000;

以上的代码执行效率很低。可以这样考虑：先将大于3000的找出来，再分组
select deptno,max(sal)from emp where sal > 3000 group by deptno;
优化策略：where和having，优先选择where，where实在完成不了，再选择 having。 
```

### where没办法解决的？

***找出每个部门的平均薪资，要求显示薪资高于2500的***

```sql
//第一步：找出每个部门的平均薪资。
select 
	deptno,avg(sal)
from	
	emp 
group by	
	deptno;
	
//第二步：要求显示薪资高于2500的
select	
	deptno,avg(sal)
from	
	emp
where		//where后面不能使用分组函数。	
	avg(sal) > 2500
group by	
	deptno;	
	
改为：优化策略。
select	
	deptno,avg(sal)
from	
	emp
group by	
	deptno
having	
	avg(sal) > 2500;
```

### 总结（单表查询结束)

```sql
select	
	...
from	
	...
where	
	...
group by	
	...
having	
	...
order by	
	...
	执行顺序：不能颠倒	
1，from 查询	2，where 过滤	3，group by 分组	4，having 接着过滤	5，select 查询	6，order by 排序输出	
从某张表中查询数据先经过where条件筛选出有价值的数据对这些有价值的数据进行分组分组之后可以使用having继续筛选经过select查询出来最后排序输出。
```

> 综合案例：
>
> 找出每个岗位的平均薪资，要求显示平均薪资大于1500的，出MANAGER之外，要求按照平均薪资降序排列。

```sql
select	job,avg(sal) as avgsalfrom	empwhere	job <> 'MANAGER'group by	jbohaving	avg(sal) > 1500order by	avgsal desc;
```

## 去除重复记录

把查询结果去除重复记录

**注意：原表数据不会被修改，只是显示结果去重。去重需要一个关键字 distinct。**

```sql
select distinct 
	job 
from 
	emp;
//select ename,distinct job from emp;
//这样编写是错误的，distinct只出现在所有字段的前方。
select distinct 
	job,deptno 
from 
	emp;
//出现在job个deptno两个字段之前，联合两个字段去重。
//统计工作岗位的数量
select 
	count(distinct job) 
from 
	emp;
```

## 连接查询

### 概述：

> 从emp表中取员工名字，从dept表中取部门名字。这种跨表查询，多张表联合起来查询数据，称为连接查询
>
> 从一张表中单独查询，称为单表查询。

### 连接查询的分类

> 根据语法的种类进行分类
>
> ​	SQL92：1992年出现的语法。
>
> ​	SQL99：1999年出现的语法。
>
> 根据表的连接方式分类
>
> ​	内连接：
>
> ​		等值连接
>
> ​		非等值连接
>
> ​		自连接
>
> ​	外连接：
>
> ​		左外连接 （左连接）
>
> ​		右外连接 （右连接）
>
> ​	全连接（不常用）

### 笛卡尔积现象

当两张表进行连接查询时，没有任何的条件限制，会发生什么现象。

> **连接查询的原理：拿出一张表的记录和另一张表的所有记录进行对比。**

```sql
//案例：查询每个员工所在的部门名称
select 
	ename, deptno 
from 
	ename,deptno;
```

> 笛卡尔积现象：当两张表进行连接查询时，没有任何的条件限制，会发生的现象是**最*终的查询结果的条数是两张表条数的乘积。***

### 怎么避免笛卡尔积现象

> 连接时接入条件，满足这个条件的记录被筛选出来。

```sql
select	
	ename,dname 
from 
	emp,denptno
where	
	emp.ename, deptno.dname;
//最终查询结果条数减少，但是执行次数不会减少。
//给表起别名，提高查询效率。
select 
	e.ename,d.dname
from 
	emp e, denptno d 
where	
	e.ename = d.dname;  
//SQL92语法。
```

***减少表的连接次数，提高程序的执行效率***

### 内连接

**特点：将能够匹配上这个查询条件的数据查询出来**

**两张表中没有主次关系**

#### 内连接之等值连接

```sql
//案例：查询每个员工所在部门的名称，显示员工和部门名
//emp e,dept d表进行连接，条件是e.deptno = d.dpetnp;
//SQL92语法
select	
	e.ename,d.dname
from	
	emp e, deptno d
where	
	e.deptno = d.deptno;	
	
//缺点：结构不清晰，表的连接条件和后期进一步筛选的条件，都放到了where的后面。

//SQL99语法
select	
	e.ename,d.dname
from	
	emp e, 
(inner) join			//inner（内部的）可以省略。	
	dept don e.deptno = d.deptno;	//条件是等值关系，所以为等值连接。
    
select 	
	...
from	
	a
(inner) join	
	bon	a和b的连接条件。
where	
	筛选条件	
	
//优点：表连接的条件是独立的，连接之后，如果还需要进一步筛选，再往后继续添加where条件即可。
```

#### 内连接之非等值连接

```sql
//案例：查询每个员工的薪资等级，要求显示员工名，薪资，员工等级
select	
	e.ename,e.sal,s.salgrade
from	
	emp e
join			//inner（内部的）可以省略。	
	salgrade son 	e.sal 
between s.local and s.hisal  
//非等值连接，条件不是一个等量关系。
```

#### 内连接之自连接

```sql
//案例：查询每个员工的上级领导，要求显示员工名和对应的领导名。
//技巧：一张表看作两张表
select	
	a.ename as '员工名',b.ename as '领导名'
from 
	emp a
join	
emp bon	a.mgr = b.empno;  //自连接的技巧，一张表看作两张表。
```

## 外连接

> 任何一个右连接都有左连接的写法
>
> 任何一个左连接都有右连接的写法
>
> 两种写法可以互换。

#### 右外连接(右连接)

```sql
select	e.ename,d.dnamefrom	emp e, right (outer) join   //outer可以省略，带着可读性强。	dept don	e.deptno = d.deptno;//right代表：表示将join关键字右边的这张表看成是主表，主要是为了将这张表的数据全部查询出来，捎带着关联查询左边的表。//在外连接中产生了主次关系
```

#### 左外连接(左连接)

```sql
select	e.ename,d.dnamefrom	emp e, left (outer) join    //outer可以省略，带着可读性强。	dept don	e.deptno = d.deptno;
```

> 外连接的查询结果的结果条数一定是 >= 内连接的查询结果的结果条数

```sql
//案例：查询每个员工的上级领导，要求显示所有的员工名字和领导名字。//外连接select	a.ename as '员工名', b.ename as ‘领导名’from	emp aleft join	emp bon 	a.mgr = b.empno;
```

### 三张表或者四张表的连接方式

```sql
//语法：select	...from	ajoin	bon	a和b的连接条件join	con	a和c的连接条件right join	don	a和d的连接条件	//一条SQL语句中，内连接和外连接可以混合出现。
```

## 子查询

### 概述：

> select语句中嵌套select语句。
>
> ```sql
> //可以出现在以下位置select	..(select).from	..(select).where	..(select).
> ```

### where子句中的子查询

```sql
//案例：找出比最低工资高的员工姓名和工资select	ename,salfrom	empwhere	sal > min(sal);//此方法不能实现，因为where子句中不能出现分组函数//实现思路：	第一步;查询最低工资	select min(sal) from emp;	第二步：	找出 > 800的	第三步：合并	select 		ename,sal 	from 		emp 	where 		sal > (select min(sal) from emp);
```

### from子句中的子查询

**注意：from后面的子查询可以将子查询的结果当作一个临时表（技巧）**

```sql
//案例：找出每个岗位的平均工资的薪资等级第一步：找出每个岗位的平均工资（按照岗位分组求平均值）select job,avg(sal) from emp group by;第二步：将第一步查询的表当作一张真实的表t，将t和salgrade进行表连接，条件：t表abg(sal) between s.losal and s.hisal;select	t.*,s.gradefrom	(select job,avg(sal) as avgsal from emp group by job) tjoin	salgrade son	t.avgsal between s.losal and s.hisal'
```

### select后面的子查询（了解即可，不常用）

**注意：这样的子查询结果不能多于1条，不然出错。**

## union的用法

> ***合并查询结果集***
>
> union的使用效率要高一些，对于表连接来说，每连接一次新表，则匹配的次数满足笛卡尔积，成倍翻。。。
>
> 但是union的效率可以减少匹配的次数。在减少匹配次数情况下，还可以完成两个结果集的拼接。
>
> a连接b连接c
>
> a 10条记录
>
> b 10条记录
>
> c 10条记录
>
> 匹配记录是：10 \* 10 \* 10 = 1000
>
> a 连接 b 一个结果： 10 \* 10 = 100;
>
> a 连接 c 一个结果： 10 \* 10 = 100;
>
> 使用union是 100 + 100 = 200次。（union把乘法变成了加法运算）

``` sql
//案例：查询工作岗位是MANAGER和SALESMAN的员工	select ename,job from emp where job = 'MANAGER' or job = 'SALESMAN';		select ename,job from emp where job in('MANAGER','SALESMAN');	union的使用：select ename,job from emp where job = 'MANAGER'unionselect ename,job from emp where job = 'SALESMAN';
```

> union的使用注意事项：
>
> union在进行结果集合并的时候，要求两个结果集的列数相同。
>
> 要求列数的数据类型也先相同。

## limit 

### 概述：

将查询结果的一部分取出，通常在分页查询中。分页的作用是为了提高用户的体验，因为一次查询出来所有的结果，用户体验差。可以翻页看。

### limit的用法：

```sql
//案例：按照薪资降序，取出排名在前五名的员工？select 	ename,salfrom	emporder by	sal desclimit 5;	//去前五//limit的语法：完整用法：limit startIndex,length	startIndex是起始下标，length是长度。缺省用法：limit 5;
```

***注意：limit是在order by之后执行的。***

```sql
//案例：取出工资排名在[3-5]名的员工。select	ename,salfrom	emporder by	sal desclimit	2, 3;//2表示起始位置从下标2开始，就是第三条记录。//3表示长度。
```

### 通用分页

> 每页显示3条记录
>
> 第一页：limit 0，3 [0 1 2]
>
> 第二页：limit 3，3 [3 4 5]
>
> 第三页：limt 6，3 [6 7 9] 
>
> 第四页：limit 9，3 [9 10 11]
>
> 每页显示pageSize条记录
>
> 第pageNo页：limit  （pageNo - 1 * pageSize)，pageSize

## 关于DQL语句的总结：

```sql
select	...from	...where	...group by	...having	...order by	...limit	...
```

## 表的创建及数据类型

### 创建表

```sql
建表的语法格式：（建表属于DDL语句，DDL包括：create drop alter）
create table 表名(字段名1 数据类型，字段名2 数据类型，字段名3 数据类型);
```

表名：建议以t\_或者tbl\_开始，可读性强。见名知义

字段名：见名知义

表名和字段名都属于标识符

### 数据类型

掌握一些常见的数据类型即可：

> varchar（长度255）：
>
> ​	可变长度的字符串，比较智能，根据实际的字符串动态	分配空间。
>
> ​	优点：节省空间
>
> ​	缺点：需要动态分配空间，速度慢
>
> char（长度255）：
>
> ​	定长字符串，不管实际的字符串长度是多少，只分配固定长度的空间取存储数据。
>
> ​	优点：不需要动态分配空间，速度快
>
> ​	缺点：使用不恰当会导致空间的浪费。
>
> ​	**varchar和char的选择：性别字段：char。姓名字段：varchar。因为性别字段长度是固定的。**
>
> int（长度11）：数字中的整型，等同于java的int
>
> bigint：数字中的长整型，相当于java中的long
>
> double：双精度浮点型数据
>
> float：单精度浮点型数据
>
> date：短日期类型
>
> datetime：长日期类型
>
> clob：字符大对象
>
> ​	最多可以存储4G的文章。
>
> ​	比如：存储一篇文章，存储一个说明。
>
> ​	超过255个字符都用CLOB字符存储。
>
> ​	Character Large OBject: CLOB
>
> blob：
>
> ​	二进制大对象
>
> ​	Binary Large OBject：BLOB
>
> ​	专门用来存储图片，声音，视频等流媒体数据。
>
> ​	往blob类型的字段上插入数据的时候，需要使用IO流。

### 创建一个学生表

```sql
create table t_student(	
    no int,    
    sex char(1),    
    name varchar(32),    
    age int(3),    
    email varchar(255)
);
```

### 快速创建表（了解）

`create table emp2 as select * from emp;`

> 原理：可以完成表的快速复制。
>
> ​		把一个查询结果当成一张表新建。
>
> ​		同时表中的数据也会存在。

### 删除表

如果这张表不存在会报错

drop table t_student;

如果这张表存在则删除。不会报错。

drop table t_student if exists;

### 对表结构的增删改

```sql
第一：在确定了表的结构之后，一般不会轻易的再次改动表，因为修改表的结构的代价比较大。修改表结构的时候，相对应的java代码也需要修改，成本是比较高的。
第二：因为成本原因，这部分的责任由设计人员承担。不需要掌握。可以使用navicat工具修改。

增加列
 alter table test add context_time datetime
 //向text表中加入context_time的列名，类型为datetime
 
```

### 插入数据，insert关键字。（DML）

### 插入值

```sql
//语法格式：	
insert into 表名（字段名1，字段名2，字段名3……）values（值1，值2，值3）;

//注意：字段名和值要一一对应，即数据类型要对应，数量要对应。
insert into t_student(no,name,sex,age,email) values (1,'zhangsan','m',20,'zhangsan@123.com');

insert into t_student(no) values (3);
//没有给其他字段指定值的时候，其他位置为NULL；
//insert语句只要执行成功了，必须会多一条记录。

//给某个字段指定默认值
create table t_student(	no int, sex char(1) defult 'f',	name varchar(32), age int(3), email varchar(255));
//sex没有插入值时，默认值是‘f’。    
//insert语句中的字段名可以省略吗，可以insert into t_student values(和字段的数据类型，数量匹配的值);
//不写字段时，相当于全写。
```

### insert一次插入多条记录

```sql
insert into t_user(id,name,birth, create_time) values (1,'zhangsan','1980-01-01'), (1,'lisi','2000-01-01')

//语法格式：
insert into t_user(字段名1，字段名2) values(),(),();
```

### 将查询结果插入到一张表当中（了解）

`creat table dept_bak as select from dept;`

`insert into dept_bak select * from dept_bak;`

很少用这种方式。

### insert插入日期

#### str_to_date函数:

**将字符串varchar类型转换成date类型**

> 语法格式：
>
> str_to_date('字符串日期'，‘日期格式’);

date_format:将date类型转换成具有一定格式的varchar的字符串类型。

***注意：数据库中的一条命名规范，所有的标识符都是全小写，单词和单词之间使用下划线分隔开。***

mysql的日期格式：

> %Y 年
>
> %m 月
>
> %d 日
>
> %h 时
>
> %i 分
>
> % 秒

```sql
insert into t_user(id,name,birth) values(1,'zhangsan',str_to_date('01-10-1990', '%d-%m-%Y'));
```

str_to_date函数可以把字符串varchar转换成date类型日期数据，

通常使用在插入insert方面，因为插入的时候需要一个日期类型的数据。需要通过函数将字符串转换成date。

如果提供的日期字符串是这个格式，str_to_date函数就不需要了，会自动做类型转换。

> %Y-%m-%d 		1990-10-01

#### date_format函数。

**将日期类型转换成特定格式的字符串**

```sql
select id name,date_format(birth,'%m-%d') from t_user;date_fromat函数的使用语法格式：date_format(日期类型，‘日期格式’)
```

## date和datetime类型的区别：

date是短日期，只包括年月日

datetime是长日期，包括年月日时分秒信息。

mysql短日期默认格式：%Y-%m-%d

mysql长日期默认格式：%Y-%m-%d %h : %i : %s;

**获取系统当前时间的函数**

now()函数并且获取的时间带有时分秒。

## DML语句

### 修改update

语法格式：

​	`update 表名 set 字段名1=值1，字段名2=值2，字段名3=值3…where 条件;`

**注意：没有条件限制会导致左右数据全部更新**

```sql
update t_user set name = 'jack',birth = '2000-10-11' where id = 2;

//修改创建时间
update t_user set name = 'jack',birth = '2000-10-11',create_time = now(); where id = 2;

//更新所有：
update t_user set name = 'abc';
```

### 删除delete语句

````sql
语法格式：
delete from 表名 where 条件
注意：没有where条件时会删除整张表的数据。

清空所有
delete from t_user;

delete语句的删除原理：（逻辑删除）
删除表中的当前内容不会删除表中的数据空间。不会释放物理空间。
缺点：删除效率低
优点：删除之后可以回滚，恢复。
````

### 快速删除表中的数据

```sql
truncate语句：快速删除表中的数据。表格式还在
缺点：不支持回滚，
优点：快速删除表，删除表的物理空间。

用法：
truncate table dept_bak;   （这种操作方式属于DDL操作）。
```

## 约束

### 约束概述

约束对应的单词：constraint

在创建表的时候，我们可以给表的字段加上一些约束，来保证这个表中数据的完整性，有效性

**约束的作用是为了保证表中的数据有效。**

### 约束包括那些

1. 非空约束   not null
2. 唯一性约束 unique
3. 主键约束  primary key (简称PK)
4. 外键约束 foregin key （简称FK）
5. 检查约束 check  （mysql不支持，oracle支持）

### 非空约束

非空约束not null约束的字段不能为null

**not null只有列级约束，没有表级约束。**

```sql
create table t_vip(
	id int,
    name varchar(255) not null
    //限定了name字段不能为空。
);
```

### SQL脚本文件

> XXXX.sql 被称为sql脚本文件，SQL脚本文件中编写了大量的sql语句，
>
> 当我们执行sql脚本文件的时候，该文件中的所有sql语句都会执行。
>
> 当需要批量执行SQL语句，可以使用SQL脚本文件
>
> 执行SQL脚本文件的方式
>
> `source + '绝对路径';`

***source命令直接执行sql脚本文件。***

### 唯一性约束，unique

唯一性约束的字段不能重复，但是可以为null。

```sql
drop table if exists t_vip;
create table t_vip(
	id int,
    name varchar(255) unique,
    //限定了name字段不能重复。
    email varchar(255)
);

insert into t_vip(id,name,email) values(1,'zhangsan','zhangsan@123.com');
insert into t_vip(idmname,email) values (2,'lisi','lisi@123.com');

insert into t_vip(id,name,email) values(1,'zhangsan','zhangsan@qq.com');
//报错：Duplicate entry ‘zhangsan’ for key ‘name’

//虽然naem被unique约束，但是可以同时插入NULL
insert into t_vip(id) values(4);
insert into t_vip(id) values(5);
```

### 两个字段联合唯一

```sql
//需求：name和email两个字段联合起来据哦于唯一性
drop table if exists t_vip;
create table t_vip(
	id int,
    name varchar(255) unique,
    email varchar(255) unique
);
//这样编写代码不符合以上唯一性需求。
//这样表示：name具有唯一性，email也具有唯一性。各自唯一。

//希望实现以下的效果：使用以上方式创建表，插入不了以下数据，因为‘zhangsan’重复了
insert into t_vip(id,name,email) values(1,'zhangsan','zhangsan@123.com');
insert into t_vip(id,name,email) values(1,'zhangsan','zhangsan@qq.com');

//使用以下方式创建
drop table if exists t_vip;
create table t_vip(
	id int,
    name varchar(255),
    email varchar(255),
    unique(name,email)
    //name和email两个字段联合起来唯一。
);
```

**约束没有添加在列的后面，这种约束被称为表级约束。**

> 需要给多个字段联合起来添加某一个约束的时候，需要使用表级约束

**约束条件添加到列后面的，叫做列级约束。**

### unique和not null可以联合吗

```sql
drop table if exists t_vip;
create table t_vip(
	id int,
    name varchar(255) not null unique
);

//被 not null unique 同时约束的话，该字段会自动变成主键字段。（注意：oracle中不一样！）

insert into t_vip(id,name) values (1,'zhangsan');
insert into t_vip(id,name) values (2,'zhangsan');
//错误：name'zhangsan'重复。

insert into t_vip(id) values(2);
//错误：name不能为空。
```

### 主键约束（primary key 简称PK *********）：

#### 主键约束的相关术语

主键约束：就是一种约束

主键字段：该字段上添加了主键约束，这样的字段叫做：主键字段

主键值：主键字段中的每一个值都叫做：主键值。

#### 主键是什么，主键的作用

主键值是每一行记录的唯一标识。

主键值是每一行记录的身份证号。

**任何一张表都应该有主键。没有主键，表无效。**

**主键不能有default**

> create table user(
>
> ​	id int(20) not null primary key ,
>
> name varchar(30) not null
>
> )engine=innodb default charset =utf8;

**主键的特征：not null + unique（主键值不能是NULL，同时也不能重复！）**

#### 给一张表添加主键：

```sql
drop table if exists t_vip;
create table t_vip(
	id int not null unique,
    name varchar(255) 
);
insert into t_vip(id ,name) values (1,'zhangsan');
insert into t_vip(id ,name) values (1,'lisi');

insert into t_vip(id ,name) values (2,'wangwu');
//错误：duplicate entry '2' for key 'PRIMARY'

insert into t_vip(name) values ('wangwu');
//错误：不能为NULL

//使用表级约束
drop table if exists t_vip;
create table t_vip(
	id int,
    name varchar(255),
    primary key(id)		//表级约束
    //一个字段做主键，叫做单一主键
);

//表级约束主要是给多个字段联合起来添加约束。
drop table if exists t_vip;
create table t_vip(
	id int,
    name varchar(255),
    email varchar(255),
    primary key(id,name)		//表级约束
    //两个字段联合起来做主键，叫做复合主键！
    //在实际开发中不建议使用复合主键。
);
```

#### 主键数量

一个表中只能添加一个主键约束。

#### 主键值建议使用类型

int

bigint

char

……

不建议使用：varchar来做主键，主键值一般都是数字，而且都是定长的。

#### 主键分类

除了单一主键和复合主键之外，还可以分为自然主键和业务主键

自然主键：主键值是一个自然数，和业务没有关系

业务主键：主键值和业务紧密联系，例如拿银行卡账号做主键值，这就是业务主键！

在实际开发中：自然主键使用比较多，因为逐渐只要做到不重复就行，不需要有意义。

业务主键不好，因为主键一旦和业务挂钩，那么当业务发生变动时，可能会影响到主键值，所以不建议业务逐渐。尽量使用自然主键。

在musql当中，有一种机制，可以帮我们自动维护主键值

```sql
drop table if exists t_vip;
create table t_vip(
	id int primary key auto_increment,
    //auto_increment从一开始，以一自增.
    name varchar(255) 
);

insert into t_vip(name) values('zhangsan');
insert into t_vip(name) values('zhangsan');
insert into t_vip(name) values('zhangsan');
insert into t_vip(name) values('zhangsan');
insert into t_vip(name) values('zhangsan');
insert into t_vip(name) values('zhangsan');
insert into t_vip(name) values('zhangsan');

select * from t_vip;
id会自增，
```

### 外键约束（foreign key 简称 FK *********）

#### 外键约束的相关术语

外键约束：就是一种约束foreign key）

外键字段：该字段上添加了外键约束，这样的字段叫做：外键字段

外键值：外键字段中的每一个值都叫做：外键值。

#### 业务背景

设计一张数据库表

第二种方案：班级一张表，学上一张表。

```sql
t_class班级表
calssno(PK)		classname
------------------------------
100				甘肃省镇原县三岔镇高三3班
101				甘肃省镇原县三岔镇高三6班

t_students
no(PK)		name		cno(FK引用_class这张表中的classno)
-----------------------------
1			jack		100
2			tom			100
3			lucy		100
4			jim			101
5			lisi		101
6			zhangsan	101
```

当cno字段没有任何约束的时候，可能会导致数据无效，可能出现一个102，但是103班级不存在，所以为了保证cno字段中的值都是100和101，需要给cno字段添加外键约束

那么：cno字段就是外键字段。cno字段中的每一个值都是外简值。

**注意：t_class 是父表。 t_students是子表。**

**删除表的顺序：先删子，再删父**

**创建表的顺序：先创建父，再创建子**

**删除数据的顺序：先删子，再删父**

**插入数据的顺序：先插入父，再插入子**

```sql
创建这张表的sql文件
set foreign_key_checks=0;
drop table if exists t_student;
drop table if exists t_class;
set foreign_key_checks=1;

create table t_class(
	classno int primary key,
    classname varchar(255)
);
create table t_student(
	no int primary key auto_increment,
    name varchar(255),
    cno int,
    constraint cforeign key(cno) references t_class(classno)
);

insert into t_class(classno,classname) values(100,'甘肃省镇原县三岔镇高三3班');
insert into t_class(classno,classname) values(101,'甘肃省镇原县三岔镇高三6班');

insert into t_student(name,cno) values('jack',100);
insert into t_student(name,cno) values('tom',100);
insert into t_student(name,cno) values('lucy',100);
insert into t_student(name,cno) values('jim',101);
insert into t_student(name,cno) values('lisi',101);
insert into t_student(name,cno) values('zhangsan',101);

select * from t_student;
select * from t_class;
```

***创建数据库时要设置数据库的编码为utf-8编码。***

外键值可以为NULL。

子表中的外键引用父表中的字段时，父表中的字段必须要主键吗？

不一定

## 存储引擎（了解内容）

### 存储引擎的作用

存储引擎时mysql中特有的一个术语，其他数据库没有，（Oracle中有，但是不叫这个名字。）

实际上是一个表存储/组织数据的方式不同

不同的存储引擎，表存储数据的方式不同。

### 给表指定存储引擎的方式

可以在建表的时候给表指定存储引擎

在见建表的时候可以在最后小括号的“）”的后面使用

ENGINE来指定存储引擎

CHARSET来指定这张表的字符编码方式

结论：

> mysql默认的存储引擎是：InnoDB
>
> mysql默认的字符编码方式是：utf-8

```sql
create table t_product(
	id int primary key,
    name varchar(255)
)engine = InnoDB default charset=utf8;
```

### mysql支持的存储引擎

`show engines \G;`

一共支持9大存储引擎。版本不同支持情况不同。

### mysql常用的存储引擎

#### MyISAM存储引擎

> 它管理的表具有以下特征
>
> ​	使用以下三个文件表示每个表
>
> ​		格式文件 - 存储表结构的定义（mytable.frm）
>
> ​		数据文件 - 存储表行的内容（nytable.MYD）
>
> ​		索引文件 - 存储表上的索引（mytable.MYI）
>
> ​	优势：可被转换为压缩，只读表来节省空间。

#### InnoDB存储引擎

> 这是mysql默认的存储引擎，同时也是重量级的存储引擎。
>
> **它支持事务，支持数据库崩溃之后的自动恢复机制**
>
> **最重要的特点是：非常安全。**
>
> 效率不是很高，不能压缩，不能节省存储空间。
>
> **它管理的表具有以下特征：**
>
> 1:支持事务(要么全成功、要么全失败);
>
> 2:行级锁定(更新数据时一般指锁定当前行)：通过索引实现、全表扫描忍让时表锁、注意间隙所的影响;
> 3:读写阻塞与事务的隔离级别相关；
> 4:具有非常高的缓存特性(既能缓存索引、也能缓存数据)；
> 5:这个表和主键以组(Cluster)的方式存储、组成一颗平衡树；
> 6:所有的辅助索引(secondary indexes)都会保存主键信息；
> 7:支持分区、表空间（逻辑名称）tablespace类似与oracle 数据库；
> 8:支持外键约束、不支持全文检索(5.5.5之前的MyISAM支持全文检索、5.5.5之后就不在支持);
> 9:相对MyISAM而言、对硬件的要求比较高

#### MEMORY存储引擎

>使用MEMORY存储引擎存储的表，其数据存储在内存中，且行的长度固定。
>
>查询效率最高，不需要和硬盘交互。
>
>但是不安全。关机之后数据消失，因为数据和索引都是存储在内存当中。

## 事务（理解，掌握\*\*\*\*\*）

### 概述：

一个事务就是一个完整的业务逻辑。

完整的业务逻辑

​	假设转账：A转账1000给B

​	A账户减去1000；

​	B账户加上1000；

​	以上的操作就是一个最小的工作单元，要么同时成功，要么同时失败，不可再分。同时成功同时失败才能保证钱是安全的。

### 只有DML语句才会有事务这个说法，其他语句没有！

insert 

update

delete

只有以上三个语句和事务有关，其他都没有关系、

因为只有以上的三个语句是数据库表中数据进行增删改查的。

只要你的操作涉及到了数据的增删改查，就一定要考虑，数据安全问题。数据安全第一位。

### 只要一条语句就能完成，还有必要存在事务机制吗？

因为所有的业务不能一条语句完成所以才需要事务的存在。

事务的本质就是要多条DML语句同时成功或者同时失败。

### 事务的实现方式

InnoDB存储引擎：提供一组用来记录事务性活动的日志文件

在事务的执行过程中，每一条DML语句的操作都会记录到‘事务性活动的日志文件’中。在事务的执行过程中，我们可以提交事务，也可以回滚事务

提交事务：

​	清空事务性活动的日志文件，将数据全部彻底持久化到数据库表中。

​	提交事务之际，事务的结束，并且是一种全部成功的结束。

回滚事务：

​	将之前所有的DML操作全部撤销，并且清空事务性活动的日志文件。

​	回滚事务标志着，事物的结束，并且是一种全部失败的结束。

### 提交事务和回滚事务

transaction：事务。

> 提交事务：commit; 语句
>
> 回滚事务：rollback; 语句 （永远只能回滚到最近的一个提交点）

mysql默认的事务提交机制：

​	mysql默认情况下是支持自动提交事务的。

​		这种的自动提交机制不符合开发习惯，日常使用中我们需要关闭这个功能。

​	没执行一条DML语句，则提交一次。

关闭mysql的自动提交机制

​	先执行这个命令：`start transaction(开启事务);`

### 事务特性

一件事就是一个事务。

A：原子性

​	说明事务是最小的工作单元。不可再分

C：一致性

​	所有事务要求，在同一个事务中，所有的操作必须同时成功或者同时失败。以保证数据的一致性

I：隔离性

​	事务之间具有一定的隔离。

D：持久性

​	事务最终结束的保障。事务提交，就相当于将没有保存到硬盘上的数据保存到硬盘上。

### 事务的4个隔离级别：

> 读未提交：read uncommitted （最低的隔离级别）
>
> ​	事务A可以读取到事务B未提交的数据。
>
> ​	这种隔离级别存在的问题是：
>
> ​		脏读现象！（Dirty Read）
>
> ​		读到了脏数据。
>
> ​	这种隔离级别是理念上的，大多数的数据库都是二档起步。
>
> 读已提交：read commited
>
> ​	事务A之能读取到事务B提交之后的数据
>
> ​	这种隔离级别解决了脏读现象。
>
> ​	这种隔离级别存在的问题是：
>
> ​		不可重复读取数据。
>
> ​		**在事务开启的时候，第一次读到的数据是3条，第二次是4条。3不等于4.称为不可重复读取。**
>
> ​	这种隔离级别是比较真实的数据，每一次读到的数据都是绝对真实。
>
> ​	oracle数据库的默认级别是read commited。
>
> 可重复读：repeatable
>
> ​	事务A开启之后，不管是多久，每一次在事务A中读取到的数据都是一致的。即使事务B将数据修改了，已提交了，事务A读取到的数据还是没有改变。
>
> ​	这种隔离级别解决了不可重复读取的问题。
>
> ​	这种隔离级别存在的问题是：
>
> ​		可能会出幻影读
>
> ​		每一次读取到的数据都是幻象，不够真实
>
> ​	mysql中默认的事务隔离级别是这个！
>
> 序列化/串行化：serializable（ 最高的隔离级别）
>
> ​	这是最高的隔离级别，效率最低，解决了所有问题。

### 演示隔离级别

查看隔离级别：select @@tx_isloation;

修改事务的隔离级别：set global transaction isolation read uncommitted;

## 索引（index）

### 概述：

索引实在数据库表的字段上添加的，是为了提高查询效率存在的一种机制。

一张表的一个字段可以添加一个索引，多个字段联合起来也可以添加索引。

索引相当于一本书的目录，是为了缩小扫描范围而存在的一种机制。

MySQL在查询方面主要就是两种方式：

​	第一种：全表扫描

​	第二种：根据索引检索

**注意：**

> 在实际中，汉语字典前面的索引也是需要排序的，并且这个索引的排序和TreeSet数据结构相同，TreeSet（TreeMap）底层是一个自平衡的二叉树！在mysql当中索引是一个B-Tree数据结构
>
> 
>
> 遵循左小右大原则存放，采用中序遍历方式遍历取数据。

### 索引的实现原理

```sql
t_user
id(PK)				name		每一行记录在硬盘上都有物理存储编号。
--------------------------------------------------------------
100					zhangsan		0x1111
120					lisi			0x2222
99					wnagwu			0x3333
88					zhangliu		0x4444
101					jack			0x5555
78					lucy			0x6666
130					tom				ox7777

1,在任何数据库当中主键上都会自动添加索引对象，id字段上自动有索引，因为id是PK，另外在mysql中，一个字段如果有unique约束的话，id字段上也会自动创建索引对象。
2，在任何数据库中，任何一张表的任何一条记录在硬盘存储上都有一个硬盘的物理存储编号。
3，在mysql中，索引是一个单独的对象，不同的存储引擎以不同的形式存在。在MyISAM存储引擎中，索引存储在一个.MYT文件中，在InnoDB存储引擎中索引存储在一个逻辑名称叫做tablespace当中，在MEMORY存储引擎中索引被存在内存中，不管索引存储在哪里，索引在mysql中都是以一个树的形式存在。（自平衡二叉树：B-Tree）
```



![index](C:/Users/Administrator/Desktop/index.PNG)

### 什么条件下，会考虑给字段上添加索引？

条件1：数据量庞大。（庞大的程度需要测试，硬件环境不同）

条件2：该字段经常出现在where的后面，以条件的形式存在，也就是说这个字段总是被扫描

条件3：该字段很少的DML操作。（因为DML之后吗，索引需要重新排序）



建议不要随便添加索引，因为索引也是需要维护的，太多的话反而会降低系统的性能。建议通过主键查询，建议通过unique约束的字段进行查询，效率是比较高的。

### 索引的创建和删除

```sql
  //创建索引
  create index emp_name_index on emp(ename)
  
  //给emp表中的ename字段添加索引，起名：emp_ename_index
  
  //删除索引
  drop index emp_ename_index on emp;
  //将emp表上的emp_ename_index索引对象删除
```

### 查看一个sql语句是否使用了索引进行检索。

explain select \* from emp where ename = 'KING';

explain: 解释。

### 索引失效

#### 失效的第一种情况

select \*from emp where ename like '%T';

ename上即使添加了索引，也不会走索引。

​	因为模糊匹配当中以“‘%"开头了

​	实际开发中，尽量避免以%开头。

​	这是一种优化的手段/策略

#### 失效的第二种情况

使用or的时候会失效，如果使用or那么要求or两边的条件字段都要有索引，才会有索引，如果其中一边有一个字段没有索引，那么另一个字段上的索引也会实现，所以这就是为什么不建议使用or的原因。

可以使用union！

#### 失效的第三种情况

使用复合索引的时候，没有使用左侧的列查找，索引失效。

复合索引：

​	两个或者更多的字段联合起来添加一个索引。

#### 失效的第四种情况

在where当中索引列参加了运算，索引失效。

#### 失效的第五种情况

在where当中索引列使用了函数

### 索引的分类

索引是各种数据库进行优化的重要手段，优化的时候优先考虑的是索引。

> 单一索引：一个字段上添加索引
>
> 复合索引：两个字段或者更多的字段上添加索引。
>
> 主键索引：主键上添加索引
>
> 唯一性索引：使用unique约束的字段上添加索引
>
> **注意：在唯一性比较弱的字段上添加索引用处不大**

## 视图（view）

### 概述：

站在不同的角度去看待同一份数据

### 怎么创建视图对象，怎么删除视图对象

```sql
//创建视图对象
create view dept2_view as select * from dept2;
//将查询结果作为视图

//删除
drop view dept2_view;

```

**注意：只有DQL语句才能创建以view的形式创建。**

`create view view_name as 这里的语句必须是DQL语句;`

### 视图的作用

可以面向视图对象进行增删改查，对**视图的增删改查，导致原表被修改。**

#### 视图对象在实际开发中的作用

假设有一条非常复杂的sql语句，而这条sql语句需要在不同的位置上反复使用，每一次都压重新编写，很长，很麻烦。怎么办？

​	可以把这条复杂的sql语句以视图对象的形式创建

​	在需要编写这条sql语句的位置直接使用视图对象，可以大大简化开发。

​	并且利于后期的维护，因为修改的时候也只需要修改一个位置就行，只需要修改视图所映射的sql语句。

**视图对象也是一个文件**，在数据库中是以文件形式存储在硬盘上的的。

### 增删改查：又叫做CRUD

CRUD是程序员之间沟通的术语，一般我们很少说增删改查

一般都说CRUD

> C: Create (增)
>
> R: Retrive （查：检索）
>
> U: Update （改）
>
> D: Delete （删）

## DBA（DataBasesAdministrator）的命令（不经常使用）

### 数据导出：

**注意：在windows的命令行窗口中**

​	导出SQL文件

​	`nysqldump bjpowernode>D:\bjpowernode.sql -uroot -p123456;`

​	导出指定的表

​	`mysqldump bjpowernode emp>D:\bjpowernode.sql -uroot -p123456;`

### 数据导入

​	**注意：需要先登录到mysql数据库服务器上**

​	然后创建数据库：create database bipowernode;

​	使用数据库：use bjpowernode;

​	然后初始化数据库：surce D:\bipowernode.sql;

## 数据库设计三范式

### 概述：

数据库的设计依据，教你怎么设计数据库表

### 数据库设计范式共有？

​	3个

1. 要求任何一张表必须有主键，每一个字段原子性不可再分

2. 建立在第一范式的基础上，要求所有非主键字段完全依赖主键，不要产生部分依赖。

3. 建立在第二范式的基础上，要求所有非主键字段直接依赖主键。不要产生传递依赖。

   ***声明：三范式是面试官经常问的，所以一定要熟记在心！***

   设计数据库表的时候，按照以上的范式进行，可以避免表中数据的冗余，空间的浪费。

### 第一范式

最核心，最重要的范式，所有表的设计都需要满足。

必须有主键，每一个字段原子性不可再分

### 第二范式

建立在第一范式的基础上，要求所有非主键字段完全依赖主键，不要产生部分依赖。

***多对多的设计方式：多对多，三张表，关系表两个外键！***

### 第三范式

建立在第二范式的基础上，要求所有非主键字段直接依赖主键。不要产生传递依赖。

***一对多，两张表，多的表加外键！***

## 总结表的设计

### 一对多

​	一对多，两张表，多的表加外键！

### 多对多

​	多对多，三张表，关系表两个外键！

### 一对一

​	表中的字段太多，太庞大的时候需要拆分成两张表。

​	***一对一，外键唯一！***

## 设计总结

**数据库是设计三范式是理论上的，最终的目的是为了客户需求，有时候会用冗余换取速度。因为在sql中，多张表的连接会降低效率。（笛卡尔积现象）同时会降低sql语句的编写难度！**

