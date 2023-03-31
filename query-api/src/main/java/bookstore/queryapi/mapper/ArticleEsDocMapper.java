package bookstore.queryapi.mapper;

import bookstore.queryapi.document.ArticleDoc;
import bookstore.queryapi.dto.ArticleDto;
import bookstore.queryapi.es_document.ArticleEsDoc;
import org.mapstruct.Mapper;
import org.mapstruct.factory.Mappers;

@Mapper
public interface ArticleEsDocMapper {
    ArticleEsDocMapper INSTANCE = Mappers.getMapper(ArticleEsDocMapper.class);

    //    @Mapping(source = "id", target = "mongoId")
    ArticleEsDoc toEsDoc(ArticleDoc mongoDoc);

    ArticleDto toDto(ArticleEsDoc doc);
}
