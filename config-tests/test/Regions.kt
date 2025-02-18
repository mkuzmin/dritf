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
        val awsRegions = aws.readRegions()
        assert(config.regions.distinct().sorted() == awsRegions.sorted())
    }
}
