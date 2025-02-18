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
    fun `region list in Gradle is actual`() {
        assert(schema.regions.sorted() == aws.readRegions().sorted())
    }

    @Test
    fun `missing regions`() {
        val diff = schema.regions - config.regions
        assert(diff.isEmpty()) {
            buildString {
                appendLine("Missing regions:")
                diff.forEach { appendLine("- $it") }
            }
        }
    }

    @Test
    fun `unknown regions`() {
        val diff = config.regions - schema.regions
        assert(diff.isEmpty()) {
            buildString {
                appendLine("Unknown regions:")
                diff.forEach { appendLine("- $it") }
            }
        }
    }
}
