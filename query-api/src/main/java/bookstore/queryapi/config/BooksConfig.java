package bookstore.queryapi.config;

import com.mongodb.MongoClientSettings;
import com.mongodb.MongoCredential;
import com.mongodb.ServerAddress;
import com.mongodb.client.MongoClient;
import com.mongodb.client.MongoClients;
import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.boot.autoconfigure.mongo.MongoProperties;
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.boot.context.properties.EnableConfigurationProperties;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.Primary;

import static java.util.Collections.singletonList;

@Configuration
//@EnableMongoRepositories(basePackageClasses = UserRepository.class, mongoTemplateRef = "primaryMongoTemplate") //or use packages
@EnableConfigurationProperties
public class BooksConfig {
    @Bean(name = "bookProps")
    @ConfigurationProperties(prefix = "spring.data.mongodb")
    @Primary
    public MongoProperties bookProps() {
        return new MongoProperties();
    }

    @Bean(name = "booksMongoClient")
    public MongoClient mongoClient(@Qualifier("bookProps") MongoProperties mongoProperties) {

        MongoCredential credential = MongoCredential
                .createCredential(mongoProperties.getUsername(), mongoProperties.getAuthenticationDatabase(), mongoProperties.getPassword());

        return MongoClients.create(MongoClientSettings.builder()
                .applyToClusterSettings(builder -> builder
                        .hosts(singletonList(new ServerAddress(mongoProperties.getHost(), mongoProperties.getPort()))))
                .credential(credential)
                .build());
    }
}
