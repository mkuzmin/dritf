@file:Suppress("unused")

import config.Config
import org.junit.jupiter.api.Test

val config = Config(directory = "../dritf.yaml")

class Regions {
    @Test
    fun `region list is sorted`() {
        assert(config.regions == config.regions.sorted())
    }

    @Test
    fun `region names are unique`() {
        val duplicates = config.regions
            .groupBy { it }
            .filter { (_, values) ->
                values.size > 1
            }
            .keys
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
