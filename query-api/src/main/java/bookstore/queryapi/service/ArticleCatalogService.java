package bookstore.queryapi.service;

import bookstore.queryapi.document.ArticleCatalogDoc;
import bookstore.queryapi.dto.ArticleCatalogDto;
import bookstore.queryapi.mapper.ArticleCatalogMapper;
import bookstore.queryapi.repo.ArticleCatalogRepo;
import org.springframework.cache.annotation.Cacheable;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.Objects;
import java.util.stream.Collectors;

@Service
public class ArticleCatalogService {
    private final ArticleCatalogRepo repo;

    public ArticleCatalogService(ArticleCatalogRepo repo) {
        this.repo = repo;
    }

    //key: query-api:article_catalogs::list
    @Cacheable(value = "article_catalogs", key = "'list'")
    public List<ArticleCatalogDto> findAllCatalogs() {
        var docs = repo.findAll();
        return buildTree(docs, null);
    }

    private List<ArticleCatalogDto> buildTree(List<ArticleCatalogDoc> docs, String parentId) {
        return docs.stream().filter(doc -> Objects.equals(parentId, doc.getParentId()))
                .map(doc -> {
                    var docDto = ArticleCatalogMapper.INSTANCE.toDto(doc);
                    docDto.setChildren(buildTree(docs, doc.getId()));
                    return docDto;
                }).collect(Collectors.toList());
    }
}
