package bookstore.queryapi.mapper;

import bookstore.queryapi.document.ArticleDoc;
import bookstore.queryapi.dto.ArticleDto;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.factory.Mappers;

@Mapper
public interface ArticleMapper {

    ArticleMapper INSTANCE = Mappers.getMapper(ArticleMapper.class);

//    @Mapping(target = "name", constant = "书籍名称测试结果")
//    @Mapping(target = "content", constant = "内容结果：书籍名称测试结果")

    @Mapping(target = "name")
    @Mapping(target = "content")
    ArticleDto toDto(ArticleDoc doc);

    ArticleDoc toDocument(ArticleDto doc);
}
