@file:Suppress("unused")

import cloudformation.Service.Companion.names
import config.Service.Companion.names
import org.junit.jupiter.api.Test

class Services {
    @Test
    fun `service list is sorted`() {
        assert(config.services.names == config.services.names.sorted())
    }

    @Test
    fun `ignored service list is sorted`() {
        assert(config.ignoredServices == config.ignoredServices.sorted())
    }

    @Test
    fun `service names are unique`() {
        val duplicates = findDuplicates(config.services.names)
        assert(duplicates.isEmpty()) { duplicates }
    }

    @Test
    fun `ignored service names are unique`() {
        val duplicates = findDuplicates(config.ignoredServices)
        assert(duplicates.isEmpty()) { duplicates }
    }

    @Test
    fun `Duplicate services`() {
        val diff = config.services.names.intersect(config.ignoredServices)
        assert(diff.isEmpty()) { "Services both listed and ignored: $diff" }
    }

    @Test
    fun `missing services`() {
        val diff = schema.services.names - config.services.names - config.ignoredServices
        assert(diff.isEmpty()) {
            buildString {
                diff.forEach { appendLine("- $it") }
            }
        }
    }

    @Test
    fun `unknown services`() {
        val diff = config.services.names + config.ignoredServices - schema.services.names
        assert(diff.isEmpty()) { diff }
    }
}
