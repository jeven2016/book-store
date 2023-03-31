package bookstore.queryapi.repo;

import bookstore.queryapi.document.ArticleDoc;
import org.bson.types.ObjectId;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.mongodb.repository.Query;
import org.springframework.data.repository.PagingAndSortingRepository;
import org.springframework.stereotype.Repository;

import java.util.Optional;

@Repository
public interface ArticleRepo extends PagingAndSortingRepository<ArticleDoc, String> {

    @Query(fields = "{'content': 0}")
    Page<ArticleDoc> findArticleDocByCatalogId(ObjectId catalogId, Pageable pageable);

    @Query(value = "{}", fields = "{'content': 0}")
    Page<ArticleDoc> findAllWithoutContent(Pageable pageable);


    @Query(value = "{ '_id': ?0}")
    Optional<ArticleDoc> findArticleDocById(String articleId);

    int countAllByCatalogId(ObjectId catalogId);
}
