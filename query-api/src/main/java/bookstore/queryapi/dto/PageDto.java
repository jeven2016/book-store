package bookstore.queryapi.dto;

import lombok.Builder;
import lombok.Getter;
import lombok.Setter;

import java.io.Serializable;
import java.util.ArrayList;
import java.util.List;

@Builder
@Getter
@Setter
public class PageDto<T extends Serializable> implements Serializable {

    private int page;

    private int limit;

    private long totalPages;

    private long count;

    @Builder.Default
    private List<T> rows = new ArrayList<>();
}
