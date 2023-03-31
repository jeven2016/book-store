package bookstore.queryapi.document;

import lombok.Getter;
import lombok.Setter;
import org.springframework.data.mongodb.core.index.Indexed;
import org.springframework.data.mongodb.core.mapping.Document;
import org.springframework.data.mongodb.core.mapping.FieldType;
import org.springframework.data.mongodb.core.mapping.MongoId;

import java.io.Serializable;
import java.time.Instant;

@Document(collection = "article")
@Getter
@Setter
public class ArticleDoc implements Serializable {
    @MongoId(FieldType.OBJECT_ID)
    private String id;

    private String name;

    private String catalogId;

    private String content;

    @Indexed
    private Instant createDate;
}
