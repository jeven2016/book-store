package bookstore.queryapi.repo;

import bookstore.queryapi.document.ArticleCatalogDoc;
import org.springframework.data.mongodb.repository.MongoRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface ArticleCatalogRepo extends MongoRepository<ArticleCatalogDoc, String> {
}
