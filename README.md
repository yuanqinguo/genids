# GENIDS

# 说明

###### 基于雪花ID算法思想，缩减时间戳41位到38位，即秒级，非毫秒级，2位作为NodeID，即只支持0,1,2, 3三个节点，11位自增ID，即支撑2^11，生成12位数字的ID值

###### 4个节点均部署后，每秒支持生成的ID数为： 2048 * 4 

###### 若想每秒生成更多的ID，有两种办法
1. 修改idworker中的位数即可，但生成的ID将会超过12位，或者直接参考标准的雪花算法生成19位的ID
2. 修改genids.yml文件pre_gen参数为true，此配置表示，会在无ID生成请求的过程中，内部协程会预先生成好50W个ID，后续的请求会从生成好的ID池中获取
   内部协程会充分利用无ID生成请求的时间，尽可能的保证永远有50W个ID可用。此方法仅能解决突然业务增加导致的ID数量剧增的问题，但若每秒都需要几万的ID时，建议使用1方法进行扩展

# 接口

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

