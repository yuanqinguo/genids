# GENIDS

# 说明

###### 基于雪花ID算法思想，缩减时间戳41位到38位，即秒级，非毫秒级，2位作为NodeID，即只支持0,1,2三个节点，11位自增ID，即支撑2^11，生成12位数字的ID值

###### 三个节点均部署后，每秒支持生成的ID数为： 2048 * 3 = 6144

###### 若想每秒生成更多的ID，修改idworker中的位数即可，但生成的ID将会超过12位，或者直接参考标准的雪花算法生成19位的ID

# 接口

## 初始化接口

```shell
/genids/init_node

query: nodeid=[0,1,2]  nodeid只能为0,1,2中的某一个，且全局唯一

curl  http://localhost:8090/genids/init_node?nodeid=0

返回：
{
    "errcode": 100, # 100表示成功，其它为失败
    "errmsg": "OK",
    "data": ""
}
```



## ID生成接口

```shell
/genids/getid

curl  http://localhost:8090/sp/genids/getid

返回：
{
    "errcode": 100, # 100表示成功，其它表示失败
    "errmsg": "OK",
    "data": {
        "id": 503105454080
    }
}
```

