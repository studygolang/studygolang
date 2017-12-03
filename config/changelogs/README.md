### 说明
数据库变更采用 `Liquibase` 管理，建表、修改字段、添加索引等操作需要编写 xml 配置文件来实现，不再需要手动改动数据库。

Liquibase 文档:   http://www.liquibase.org/documentation/index.html

`changelogs` 文件夹下建议按照 `1.0、1.1、1.2、2.0 ...` 等存放每次需要改动的配置文件。

每个版本中的 `xml` 文件名需要和 `env.ini` 中配置 `dbname` 的一致，比如 `studygolang.xml`.

### 示例

`changelogs/1.0/studygolang.xml` 中新建了一个表 `test_liquibase`.

```xml
<changeSet id="1" author="javasgl">
    <createTable tableName="test_liquibase">
        <column name="id" type="int">
            <constraints primaryKey="true" nullable="false"/>
        </column>
        <column name="name" type="varchar(50)">
            <constraints nullable="false"/>
        </column>
        <column name="active" type="boolean" defaultValueBoolean="true"/>
    </createTable>
</changeSet>
```

执行

```
./bin/migrator --changeVersion=1.0
```

即可在数据中新建一个表 `test_liquibase`。

---

过了一段时间，需要给这个表添加一个字段 `status`。编写 `xml` 配置文件存于 `changelogs/1.1/studygolang.xml`，内容如下：

```xml
<changeSet id="1" author="javasgl">
    <comment>增加 status 字段</comment>
    <addColumn tableName="test_liquibase">
        <column name="status" type="tinyint unsigned" defaultValue="0" remarks="status">
            <constraints nullable="false" />
        </column>
    </addColumn>
</changeSet>
```

执行:

```
./bin/migrator --changeVersion=1.1
```

即可为 `test_liquibase` 表加上 `status` 字段。

> Liquibase 更多功能请看其 [官方文档](http://www.liquibase.org/documentation/index.html)，功能很强大。