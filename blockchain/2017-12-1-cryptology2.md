---
layout:     post
title:      "《图解密码技术》笔记（二）"
subtitle:   "《图解密码技术》笔记"
date:       2017-12-1
author:     "yucs"
catalog:    true
categories: 
	- 区块链

tags:
    - 区块链
      
---
# 单向散列函数
- 可以保证消息的完整性，但不能对消息进行认证 （可以辨别出篡改，无法辨别出伪装）
- 散列值得长度是固定的，和消息长度无关。
![单向散列函数](https://yucs.github.io/picture/单向散列函数.png)
- 性质：
  - 根据任意长度的消息计算出固定长度的散列值。
  - 能够快速计算出散列值。  
  - 如果两个消息体产生同个散列值称为碰撞,密码技术中所使用的单向散列函数必须具备抗碰撞性。
  - **哪怕只有1比特的改变，也必须有很高的概率产生不同的散列值**
  - ![散列值碰撞](https://yucs.github.io/picture/散列值碰撞.png)
  - 具备单向性：
    ![散列值单向性](https://yucs.github.io/picture/散列值单向性.png)
- 实际应用：
  - MD5来验证是否同个文件软件。
  - 数字签名处理过程非常耗时，因此一般不会对整个消息内容直接进行数字签名，而是先通过单向散列函数计算出消息的散列值，然后再对整个散列值进行数字签名。
 - SHA-1的抗碰撞性已被攻破。
 - SHA-256,SHA-384,SHA-512这些统称SHA-2，尚未被攻破。
 - SHA3简介：由于近年来对传统常用Hash 函数如MD4、MD5、SHA0、SHA1、RIPENMD 等的成功攻击2012年10月2日，Keccak被选为NIST竞赛的胜利者， 成为SHA-3。


# 数字签名
 - 公钥密码和数字签名：
   ![密钥使用方式](https://yucs.github.io/picture/密钥使用方式.png)
  ![密钥加密](https://yucs.github.io/picture/密钥加密.png)
  ![数字签名](https://yucs.github.io/picture/数字签名.png)

 -  数字签名是对消息的散列值签名：
   ![数字签名流程](https://yucs.github.io/picture/数字签名流程.png) 
   ![数字签名流程2](https://yucs.github.io/picture/数字签名流程2.png) 
 - RAS的数字签名和验证：
  ![RAS数字签名](https://yucs.github.io/picture/RAS数字签名.png)
 - 数字签名不能保证机密性，在数字签名中，只有发送者才持有生成签名的私钥，防止否认。
 - 对称密码的密钥是机密性的精华，单向散列函数的散列值是完整性的精华。
 - 数字签名是非常重要的认证技术，但前提是用于验证签名的发送者的公钥没有被伪造，即要确认公钥是否合法，可以对公钥施加数字签名，之就是证书。

# 证书
公钥证书（public-Key Certificate PKC）由可信的认证机构（certification Authority ,CA）施加数字签名。公钥证书也简称证书。

- 认证例子：
![证书认证](https://yucs.github.io/picture/证书认证.png)

- 证书标准规范X.509:
 ![X.509](https://yucs.github.io/picture/X.509.png)

简单例子：

```sh
#生成证书
openssl genrsa -out ca/ca-key.pem 1024  
#查看证书信息
openssl x509 -in cert.pem -noout -text
```

- 公钥基础设施（PKI）
  PKI(public-key infrastructure)是为了能够有效地运用公钥而制定的一系列规范和规格的总称。
  - 组成要素： 用户：使用PKI的人，认证机构（CA）：颁发证书的人 ， 仓库：保存证书的数据库：
  ![pki](https://yucs.github.io/picture/PKI.png)   
 
 - CA的工作：
   ![pki2](https://yucs.github.io/picture/PKI2.png)

 - CRL(Certificate Revocation list): 证书作废清单，当用户的私钥丢失，被盗是，认证机构需要对证书进行作废。

- 认证机构的层次：
  ![pki3](https://yucs.github.io/picture/PKI3.png)

- 其他参考：
  - [数字证书原理](http://www.cnblogs.com/JeffreySun/archive/2010/06/24/1627247.html)：文中首先解释了加密解密的一些基础知识和概念，然后通过一个加密通信过程的例子说明了加密算法的作用，以及数字证书的出现所起的作用。
 - [那些证书相关的玩意儿:X.509,PEM,DER,CRT,CER,KEY,CSR,P12](http://www.cnblogs.com/guogangj/p/4118605.html)
     - CSR（Certificate Signing Request）：即证书签名请求,这个并不是证书,而是向权威证书颁发机构获得签名证书的申请。
     - PEM（Privacy Enhanced Mail）,打开看文本格式,以"-----BEGIN..."开头, "-----END..."结尾,内容是BASE64编码.
 


# 随机数-不可预测性的源泉
- 随机数用来生成对称密钥，公钥密钥。
- 性质
  -  随机性：不存在统计学偏差，完全杂乱的数列
  -  不可预测性：不能从过去的数列推测出下一个出现的数
  -  不可重复性：除非将数列本身保存下来，否则不能重现相同的数列.
  - ![随机数](https://yucs.github.io/picture/随机数.png)
  - 伪随机数生成器
   - 根据外部输入的种子生产伪随机数列
   ![伪随机生成器](https://yucs.github.io/picture/伪随机生成器.png)
   - 伪随机数生成器是公开，但种子需要自己保密，这类似密码算法是公开，而密钥要自己保密。
   - 具体的伪随机数生成器
    -  很多伪随机数生成器的库函数使用线性同余法编写，但是不具备不可预测性，不能用于密码技术。
   ![线性同余法](https://yucs.github.io/picture/线性同余法.png) 
    - 用密码实现伪随机数生成器
  ![密码实现伪随机生成器](https://yucs.github.io/picture/密码实现伪随机生成器.png) 

# 小结
 ![密码工具箱小结](https://yucs.github.io/picture/密码工具箱小结.png) 
- 密钥是机密性的精华
- 散列值是完整性的精华
- 种子是不可预测性的精华
-  如果量子密码计算机进入实用领域，就能产生完美的密码技术
-  如果量子计算机比量子密码先进入实用领域，则使用目前的密码技术所产生的密文将全部破译。
-  即使真的拥有完美的密码技术，也不可能实现完美的安全性，因为必须会有人类-即不完美的我们参与其中。 