import config.Config

val config = Config(directory = "../dritf.yaml")

fun findDuplicates(data: List<String>): List<String> =
    data
        .groupBy { it }
        .filter { (_, values) ->
            values.size > 1
        }
        .keys.toList()
