# prometheus-test
普罗米修斯-测试
## 说明:
- 启动服务 之前先修改配置 
  - 配置在conf文件夹
  - 启动参数 -c 可以指定配置文件
## 数据库表sql
```mysql

create table images
(
    id         int UNSIGNED auto_increment comment 'ID' primary key,
    created_at int default 0 not null,
    updated_at int default 0 not null,
    url        text          not null
);

create table keywords
(
    id         int UNSIGNED auto_increment comment 'ID' primary key,
    created_at int default 0                     not null,
    updated_at int default 0                     not null,
    content    varchar(1024) collate utf8mb4_bin not null,
    md5        char(32)                          not null
) ENGINE = MEMORY;

create index idx_md5
    on keywords (md5) USING HASH;

create table image_mappings
(
    id         int UNSIGNED auto_increment comment 'ID' primary key,
    created_at int default 0 not null,
    updated_at int default 0 not null,
    image_id   int UNSIGNED  not null,
    keyword_id int UNSIGNED  not null,
    constraint unique_image_key
        unique (image_id, keyword_id)
);

create index idx_keyword_id
    on image_mappings (keyword_id);


```