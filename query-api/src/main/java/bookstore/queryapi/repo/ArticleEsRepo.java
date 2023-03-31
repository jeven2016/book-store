package bookstore.queryapi.repo;

import bookstore.queryapi.es_document.ArticleEsDoc;
import org.springframework.data.domain.Pageable;
import org.springframework.data.elasticsearch.annotations.Highlight;
import org.springframework.data.elasticsearch.annotations.HighlightField;
import org.springframework.data.elasticsearch.annotations.HighlightParameters;
import org.springframework.data.elasticsearch.annotations.SourceFilters;
import org.springframework.data.elasticsearch.core.SearchHits;
import org.springframework.data.elasticsearch.repository.ElasticsearchRepository;
import org.springframework.stereotype.Repository;


@Repository
public interface ArticleEsRepo extends ElasticsearchRepository<ArticleEsDoc, String> {

    //Sometimes the user does not need to have all the properties of an entity returned from a search but only a subset.
    //Elasticsearch provides source filtering to reduce the amount of data that is transferred across the network to the application.
    @Highlight(
            fields = {
                    @HighlightField(name = "name", parameters = @HighlightParameters(
                            preTags = "<span style='color:red'>",
                            postTags = "</span>",
                            fragmentSize = 500,
                            numberOfFragments = 3
                    ))
            },
            parameters = @HighlightParameters(
                    preTags = "<span style='color:red'>",
                    postTags = "</span>",
                    fragmentSize = 500,
                    numberOfFragments = 3
            )
    )
    @SourceFilters(excludes = "content")
    SearchHits<ArticleEsDoc> findArticleEsDocByName(String name, Pageable pageable);

}
