import cloudformation.Schema
import config.Config

val config = Config(directory = "../dritf.yaml")
val schema = Schema(directory = "build/schema")

fun findDuplicates(data: List<String>): List<String> =
    data
        .groupBy { it }
        .filter { (_, values) ->
            values.size > 1
        }
        .keys.toList()
