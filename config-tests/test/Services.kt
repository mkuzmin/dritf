@file:Suppress("unused")

import config.Service.Companion.names
import org.junit.jupiter.api.Test

class Services {
    @Test
    fun `service list is sorted`() {
        assert(config.services.names == config.services.names.sorted())
    }

    @Test
    fun `service names are unique`() {
        val duplicates = findDuplicates(config.services.names)
        assert(duplicates.isEmpty()) {
            "Duplicate services: $duplicates"
        }
    }
}
