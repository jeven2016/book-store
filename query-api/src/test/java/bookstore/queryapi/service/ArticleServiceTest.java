package bookstore.queryapi.service;

import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;

@SpringBootTest
public class ArticleServiceTest {

    @Autowired
    ArticleService service;


    @Test
    public void testInsertEs() {
        service.insertArticlesToEs();
    }
}
