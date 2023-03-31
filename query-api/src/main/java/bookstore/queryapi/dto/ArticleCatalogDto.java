package bookstore.queryapi.dto;

import lombok.Getter;
import lombok.Setter;

import java.io.Serializable;
import java.time.Instant;
import java.util.Collections;
import java.util.List;

@Getter
@Setter
public class ArticleCatalogDto implements Serializable {
    private String id;

    private String name;
    private int order;
    private int articleCount;
    private String description;

    private Instant createDate;
    private Instant lastUpdate;

    private List<ArticleCatalogDto> children = Collections.emptyList();
}
