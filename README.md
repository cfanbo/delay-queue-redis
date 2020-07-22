# delay-queue-redis

基于redis实现的延时队列

主要用到redis中的 有序集合（sorted sets）和 散列（hashes）两种数据类型。
hashes里存储的是实际的数据信息，sets里存储的是权重，以秒为时间单位来进行数据的消费

