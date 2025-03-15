@file:Suppress("unused")

import org.junit.jupiter.api.Test
import org.junit.jupiter.api.assertAll

class Terraform {
    @Test
    fun `terraform resource type is missing`() {
        assertAll(config.services.flatMap { service ->
            service.resourceTypes.map { rt ->
                {
                    assert(rt.terraform != null) { "${service.name}/${rt.name} doesn't have terraform mapping" }
                }
            }
        })
    }

    @Test
    fun `terraform resource types are unique`() {
        val tfTypes = config.services.flatMap { service ->
            service.resourceTypes.mapNotNull { rt ->
                rt.terraform?.name
            }
        }
        val duplicates = findDuplicates(tfTypes)
        assert(duplicates.isEmpty()) {
            "Duplicate Terraform types: $duplicates"
        }
    }
}
