package bookstore.queryapi.document;

import lombok.Getter;
import lombok.Setter;
import org.springframework.data.mongodb.core.index.Indexed;
import org.springframework.data.mongodb.core.mapping.Document;
import org.springframework.data.mongodb.core.mapping.FieldType;
import org.springframework.data.mongodb.core.mapping.MongoId;

import java.io.Serializable;
import java.time.Instant;

@Document(collection = "catalog", collation = "zh")//Collation特性允许MongoDB的用户根据不同的语言定制排序规则
@Getter
@Setter
public class ArticleCatalogDoc implements Serializable {

    //@MongoId String id;  => The id is treated as String without further conversion.
    //@MongoId ObjectId id; => The id is treated as ObjectId.
    //@MongoId(FieldType.OBJECT_ID) String id; =>The id is treated as ObjectId
    // if the given String is a valid ObjectId hex, otherwise as String. Corresponds to @Id usage.
    @MongoId(FieldType.OBJECT_ID)
    private String id;

    private String parentId;

    @Indexed
    private String name;
    private int order;
    private int articleCount;
    private String description;

    @Indexed
    private Instant createDate;
    private Instant lastUpdate;

//    @DBRef
//    private List<ArticleCatalog> catalogs = new ArrayList<>();

}
