package bookstore.queryapi.dto;

import lombok.Getter;
import lombok.Setter;

import java.io.Serializable;
import java.time.Instant;

@Getter
@Setter
public class ArticleDto implements Serializable {

    private String id;

    private String name;

    private String content;

    private String catalogId;

    private Instant createDate;
}
