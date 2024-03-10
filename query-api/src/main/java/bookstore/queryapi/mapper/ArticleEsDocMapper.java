package bookstore.queryapi.mapper;

import bookstore.queryapi.document.ArticleDoc;
import bookstore.queryapi.dto.ArticleDto;
import bookstore.queryapi.es_document.ArticleEsDoc;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.ReportingPolicy;
import org.mapstruct.factory.Mappers;

@Mapper(unmappedTargetPolicy = ReportingPolicy.IGNORE)
public interface ArticleEsDocMapper {
    ArticleEsDocMapper INSTANCE = Mappers.getMapper(ArticleEsDocMapper.class);

    //    @Mapping(source = "id", target = "mongoId")
    ArticleEsDoc toEsDoc(ArticleDoc mongoDoc);

//    @Mapping(target = "name", constant = "书籍名称测试结果")
//    @Mapping(target = "content", constant = "内容结果：书籍名称测试结果")

    @Mapping(target = "name")
    @Mapping(target = "content")
    ArticleDto toDto(ArticleEsDoc doc);
}
