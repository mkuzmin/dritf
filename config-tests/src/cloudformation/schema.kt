package cloudformation

data class Schema(
    val regions: List<String>,
    val services: List<Service>,
) {
    companion object {
        operator fun invoke(directory: String): Schema {
            val regions = readRegions()
            return Schema(
                regions = regions,
                services = listOf()
//                services = readServices(regions, directory)
            )
        }

        private fun readRegions(): List<String> {
            val filename = Companion::class.java.getResource("/regions.txt")
                ?: error("Could not find the file")

            val regions = filename
                .readText()
                .lines()
                .filter { it.isNotBlank() }

            return regions
        }

    }
}

data class Service(
    val name: String,
    val resourceTypes: List<ResourceType>,
) {
    companion object {
        val List<Service>.names
            get() = this.map { it.name }
    }
}

data class ResourceType(
    val name: String,
    val regions: List<String>,
    val isListable: Boolean,
    val parents: List<String>,
) {
    companion object {
        val List<ResourceType>.names
            get() = this.map { it.name }
    }
}
