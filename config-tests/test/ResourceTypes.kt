@file:Suppress("unused")

import cloudformation.ResourceType.Companion.names
import config.ResourceType.Companion.names
import org.junit.jupiter.api.Test
import org.junit.jupiter.api.assertAll

class ResourceTypes {
    @Test
    fun `types are sorted`() {
        assertAll(config.services.map { service ->
            {
                assert(service.resourceTypes.names == service.resourceTypes.names.sorted()) {
                    "service ${service.name}"
                }
            }
        })
    }

    @Test
    fun `ignored types are sorted`() {
        assertAll(config.services.map { service ->
            {
                assert(service.ignoredTypes == service.ignoredTypes.sorted()) {
                    "service ${service.name}"
                }
            }
        })
    }

    @Test
    fun `types are unique`() {
        assertAll(config.services.map { service ->
            {
                val duplicates = findDuplicates(service.resourceTypes.names)
                assert(duplicates.isEmpty()) { "service ${service.name}: $duplicates" }
            }
        })
    }

    @Test
    fun `ignored types are unique`() {
        assertAll(config.services.map { service ->
            {
                val duplicates = findDuplicates(service.ignoredTypes)
                assert(duplicates.isEmpty()) { "service ${service.name}: $duplicates" }
            }
        })
    }

    @Test
    fun `types are both listed and ignored`() {
        assertAll(config.services.map { service ->
            {
                val diff = service.resourceTypes.names.intersect(service.ignoredTypes)
                assert(diff.isEmpty()) { "service ${service.name}: $diff" }
            }
        })
    }

    @Test
    fun `missing types`() {
        assertAll(config.services.map { service ->
            {
                val schemaNames = schema.services.single { it.name == service.name }.resourceTypes.names
                val diff = schemaNames - service.resourceTypes.names - service.ignoredTypes
                assert(diff.isEmpty()) {
                    buildString {
                        appendLine("Missing resource types in service ${service.name}:")
                        diff.forEach { appendLine("      - $it") }
                    }
                }
            }
        })
    }

    @Test
    fun `unknown types`() {
        assertAll(config.services.map { service ->
            {
                val schemaNames = schema.services.single { it.name == service.name }.resourceTypes.names
                val diff = service.resourceTypes.names + service.ignoredTypes - schemaNames
                assert(diff.isEmpty()) { "service ${service.name}: $diff" }
            }
        })
    }

    @Test
    fun `regions are sorted`() {
        assertAll(config.services.flatMap { service ->
            service.resourceTypes.map { rt ->
                {
                    assert(rt.regions == rt.regions.sorted()) {
                        "${service.name}/${rt.name}"
                    }
                }
            }
        })
    }

    @Test
    fun `regions are unique`() {
        assertAll(config.services.flatMap { service ->
            service.resourceTypes.map { rt ->
                {
                    val duplicates = findDuplicates(rt.regions)
                    assert(duplicates.isEmpty()) { "${service.name}/${rt.name}: $duplicates" }
                }
            }
        })
    }

    @Test
    fun `missing regions`() {
        assertAll(config.services.flatMap { service ->
            service.resourceTypes.map { rt ->
                {
                    val configRegions = if (rt.regions.isEmpty()) config.regions else rt.regions
                    val schemaRegions = schema
                        .services.single { it.name == service.name }
                        .resourceTypes.single { it.name == rt.name }
                        .regions

                    val diff = schemaRegions - configRegions
                    assert(diff.isEmpty()) { "${service.name}/${rt.name}: $diff" }
                }
            }
        })
    }

    @Test
    fun `unknown regions`() {
        assertAll(config.services.flatMap { service ->
            service.resourceTypes.map { rt ->
                {
                    val configRegions = if (rt.regions.isEmpty()) config.regions else rt.regions
                    val schemaRegions = schema
                        .services.single { it.name == service.name }
                        .resourceTypes.single { it.name == rt.name }
                        .regions

                    val diff = configRegions - schemaRegions
                    assert(diff.isEmpty()) { "${service.name}/${rt.name}: $diff" }
                }
            }
        })
    }
}
