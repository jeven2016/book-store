package bookstore.queryapi.mapper;

import bookstore.queryapi.document.ArticleDoc;
import bookstore.queryapi.dto.ArticleDto;
import org.mapstruct.Mapper;
import org.mapstruct.factory.Mappers;

@Mapper
public interface ArticleMapper {

    ArticleMapper INSTANCE = Mappers.getMapper(ArticleMapper.class);

    ArticleDto toDto(ArticleDoc doc);

    ArticleDoc toDocument(ArticleDto doc);
}
