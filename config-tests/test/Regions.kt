@file:Suppress("unused")

import org.junit.jupiter.api.Test

class Regions {
    @Test
    fun `region list is sorted`() {
        assert(config.regions == config.regions.sorted())
    }

    @Test
    fun `region names are unique`() {
        val duplicates = findDuplicates(config.regions)
        assert(duplicates.isEmpty()) {
            "Duplicate regions: $duplicates"
        }
    }

    @Test
    fun `region list is actual`() {
        assert(config.regions.distinct().sorted() == awsRegions.sorted())
    }

    @Test
    fun `region list in Gradle is actual`() {
        val filename = javaClass.getResource("regions.txt")
            ?: error("Could not find the file")
        val regions = filename
            .readText()
            .lines()
            .filter { it.isNotBlank() }

        assert(regions.sorted() == awsRegions.sorted())
    }
}
