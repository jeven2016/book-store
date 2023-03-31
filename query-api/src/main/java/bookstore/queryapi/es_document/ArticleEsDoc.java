package bookstore.queryapi.es_document;

import lombok.Getter;
import lombok.Setter;
import org.springframework.data.annotation.Id;
import org.springframework.data.annotation.Transient;
import org.springframework.data.elasticsearch.annotations.Document;
import org.springframework.data.elasticsearch.annotations.Field;
import org.springframework.data.elasticsearch.annotations.FieldType;

import java.io.Serializable;
import java.time.Instant;

@Document(indexName = "article")
@Getter
@Setter
public class ArticleEsDoc implements Serializable {

    @Id
    private String id;

    //通过repository保存，会提交正确的mapping， 通过template无法提交document上的mapping映射
    @Field(type = FieldType.Text, analyzer = "ik_max_word", searchAnalyzer = "ik_smart")
    private String name;

    @Field(type = FieldType.Keyword)
    private String catalogId;

    @Transient
//    @Field(type = FieldType.Text, index = false)
    private String content;

    @Field(type = FieldType.Date)
    private Instant createDate;
}
