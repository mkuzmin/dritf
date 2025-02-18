package cloudformation

import kotlinx.serialization.Serializable

@Serializable
data class AwsType(
    val typeName: String,
    val handlers: Handlers? = null,
)

@Serializable
data class Handlers(
    val list: Handler? = null,
)

@Serializable
data class Handler(
    val permissions: List<String>,
    val handlerSchema: HandlerSchema? = null,
)

@Serializable
data class HandlerSchema(
    val required: List<String>? = null,
)
