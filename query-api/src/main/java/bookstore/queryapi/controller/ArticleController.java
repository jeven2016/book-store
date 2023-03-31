package bookstore.queryapi.controller;

import bookstore.queryapi.dto.ArticleDto;
import bookstore.queryapi.dto.PageDto;
import bookstore.queryapi.mapper.ArticleEsDocMapper;
import bookstore.queryapi.mapper.ArticleMapper;
import bookstore.queryapi.service.ArticleService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.domain.Page;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.Objects;
import java.util.stream.Collectors;

@RestController
public class ArticleController {
    private ArticleService articleService;

    @Autowired
    public void setArticleService(ArticleService articleService) {
        this.articleService = articleService;
    }

    @GetMapping("/v1/article-catalogs/{articleCatalogId}/articles")
    public PageDto<ArticleDto> findAll(@PathVariable String articleCatalogId,
                                       @RequestParam(name = "limit", required = false) Integer limit,
                                       @RequestParam(name = "page", required = false) Integer page) {
        limit = Objects.requireNonNullElse(limit, 10);
        page = Objects.requireNonNullElse(page, 1) - 1;

        Page<ArticleDto> pageInfo = articleService.getPage(articleCatalogId, page, limit)
                .map(ArticleMapper.INSTANCE::toDto);

        return PageDto.<ArticleDto>builder().page(page).limit(limit).count(pageInfo.getTotalElements())
                .totalPages(pageInfo.getTotalPages())
                .rows(pageInfo.stream().collect(Collectors.toList())).build();
    }

    @GetMapping("/v1/articles/{articleId}")
    public ResponseEntity<ArticleDto> findById(@PathVariable String articleId) {
        var articleDoc = articleService.findById(articleId);
        return articleDoc.map(doc -> ResponseEntity.ok(ArticleMapper.INSTANCE.toDto(doc)))
                .orElseGet(() -> ResponseEntity.notFound().build());

    }

    @GetMapping("/v1/articles")
    public PageDto<ArticleDto> search(@RequestParam String name,
                                      @RequestParam String catalogId,
                                      @RequestParam(name = "limit", required = false) Integer limit,
                                      @RequestParam(name = "page", required = false) Integer page) {
        limit = Objects.requireNonNullElse(limit, 10);
        page = Objects.requireNonNullElse(page, 1) - 1;

        var pageBuilder = PageDto.<ArticleDto>builder();

        var hits = articleService.search(name, catalogId, page, limit);
        List<ArticleDto> rows = hits.getSearchHits().stream().map(searchHit -> {
            var dto = ArticleEsDocMapper.INSTANCE.toDto(searchHit.getContent());

            //高亮属性
            var hlFiledValues = searchHit.getHighlightField("name");
            if (!hlFiledValues.isEmpty()) {
                dto.setName(hlFiledValues.get(0));
            }
            return dto;
        }).toList();

        var totalPages = Math.ceil((double) (hits.getTotalHits() - 1) / limit + 1);
        pageBuilder.count(hits.getTotalHits())
                .totalPages((long) totalPages)
                .page(page).limit(limit)
                .rows(rows);

        return pageBuilder.build();
    }


    @GetMapping("/v1/articles/sync")
    @ResponseStatus(HttpStatus.ACCEPTED)
    public void insert() {
        articleService.syncToEs();
    }
}
