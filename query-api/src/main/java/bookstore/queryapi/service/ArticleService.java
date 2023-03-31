package bookstore.queryapi.service;

import bookstore.queryapi.document.ArticleDoc;
import bookstore.queryapi.es_document.ArticleEsDoc;
import bookstore.queryapi.mapper.ArticleEsDocMapper;
import bookstore.queryapi.repo.ArticleEsRepo;
import bookstore.queryapi.repo.ArticleRepo;
import lombok.extern.slf4j.Slf4j;
import org.bson.types.ObjectId;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.cache.annotation.CacheConfig;
import org.springframework.cache.annotation.Cacheable;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Sort;
import org.springframework.data.elasticsearch.client.elc.NativeQuery;
import org.springframework.data.elasticsearch.core.ElasticsearchOperations;
import org.springframework.data.elasticsearch.core.SearchHits;
import org.springframework.data.elasticsearch.core.query.HighlightQuery;
import org.springframework.data.elasticsearch.core.query.highlight.Highlight;
import org.springframework.data.elasticsearch.core.query.highlight.HighlightField;
import org.springframework.data.elasticsearch.core.query.highlight.HighlightParameters;
import org.springframework.scheduling.annotation.Async;
import org.springframework.stereotype.Service;
import org.springframework.util.StringUtils;

import java.util.List;
import java.util.Optional;
import java.util.stream.Collectors;

@Slf4j
@Service
@CacheConfig(cacheNames = "articles")
public class ArticleService {
    private ArticleRepo articleRepo;

    private ArticleEsRepo esRepo;

    private volatile boolean syncRunning = false;

    private ElasticsearchOperations searchOperations;

    @Autowired
    public void setElasticsearchOperations(ElasticsearchOperations searchOperations) {
        this.searchOperations = searchOperations;
    }

    @Autowired
    public void setEsRepo(ArticleEsRepo esRepo) {
        this.esRepo = esRepo;
    }

    @Autowired
    public void setArticleRepo(ArticleRepo articleRepo) {
        this.articleRepo = articleRepo;
    }

    @Cacheable
    public Page<ArticleDoc> getPage(String articleCatalogId, int page, int limit) {
        log.info("query a page of records from db");
        var pageRequest = PageRequest.of(page, limit, Sort.by(Sort.Direction.DESC, "createDate"));
        return articleRepo.findArticleDocByCatalogId(new ObjectId(articleCatalogId), pageRequest);
    }

    public Page<ArticleDoc> getPageWithoutCache(int page, int limit) {
        log.info("query a page of records from db");
        var pageRequest = PageRequest.of(page, limit);
        return articleRepo.findAllWithoutContent(pageRequest);
    }

    @Cacheable(key = "#articleId")
    public Optional<ArticleDoc> findById(String articleId) {
        return articleRepo.findArticleDocById(articleId);
    }

    public int count(String articleCatalogId) {
        return articleRepo.countAllByCatalogId(new ObjectId(articleCatalogId));
    }

    public SearchHits<ArticleEsDoc> search(String name, String catalogId, int page, int limit) {
        var pageRequest = PageRequest.of(page, limit);

        var hlParams = HighlightParameters.builder().withPreTags("<span style='color:red'>")
                .withPostTags("</span>")
                .withFragmentSize(500)
                .withNumberOfFragments(3).build();
        var fields = List.of(new HighlightField("name"));

        var hl = new Highlight(hlParams, fields);
        var hlQuery = new HighlightQuery(hl, ArticleEsDoc.class);
        var queryBuilder = NativeQuery.builder();

        //名称分词匹配
        if (StringUtils.hasText(name)) {
            queryBuilder.withQuery(q -> q.match(m -> m.field("name").query(name)));
        }

        // 匹配catalogId
        if (StringUtils.hasText(catalogId)) {
            queryBuilder.withQuery(q -> q.term(t -> t.field("catalogId").value(catalogId)));
        }

        //不包含content内容
        var query = queryBuilder.withFields("id", "name", "catalogId", "createDate")
                .withPageable(pageRequest)
                .withHighlightQuery(hlQuery)
                .build();

        return searchOperations.search(query, ArticleEsDoc.class);
    }


    @Async
    public void syncToEs() {
        if (syncRunning) {
            log.info("a sync task is running, please try later");
            return;
        }
        var initPage = 0;
        Page<ArticleDoc> pageInfo;
        try {
            do {
                pageInfo = this.getPageWithoutCache(initPage, 10);
                var esDocs = pageInfo.stream().map(ArticleEsDocMapper.INSTANCE::toEsDoc).collect(Collectors.toList());
                esRepo.saveAll(esDocs);

                log.info("page {}/{} is processed, {} records", initPage + 1,
                        pageInfo.getTotalPages(), pageInfo.getNumberOfElements());
                initPage++;
            } while (initPage < pageInfo.getTotalPages());
        } finally {
            syncRunning = false;
        }
    }
}
