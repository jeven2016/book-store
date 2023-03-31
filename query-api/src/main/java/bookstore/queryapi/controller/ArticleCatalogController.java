package bookstore.queryapi.controller;

import bookstore.queryapi.dto.ArticleCatalogDto;
import bookstore.queryapi.service.ArticleCatalogService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.List;

@RestController
public class ArticleCatalogController {
    private ArticleCatalogService service;


    @Autowired
    public void setService(ArticleCatalogService service) {
        this.service = service;
    }


    @GetMapping("/v1/article-catalogs")
    public List<ArticleCatalogDto> getArticleCatalogs() {
        return service.findAllCatalogs();
    }

    @GetMapping("/test")
    public String test() {
        return "test";
    }
}
