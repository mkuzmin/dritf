package config

import com.charleskorn.kaml.Yaml
import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable
import java.io.File

@Serializable
data class Config(
    val regions: List<String>,
    val services: List<Service>,
) {
    companion object {
        operator fun invoke(directory: String) =
            Yaml.default.decodeFromString(
                serializer(),
                File(directory).readText()
            )
    }
}

@Serializable
data class Service(
    val name: String,
    @SerialName("resource_types") val resourceTypes: List<ResourceType>,
)

@Serializable
data class ResourceType(
    val name: String,
)
