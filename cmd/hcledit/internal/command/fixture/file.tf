module "my-module" {
  source = "source.tar.gz"

  bool_variable   = true
  int_variable    = 1
  string_variable = "string"

  # Comment
  array_variable = ["a", "b", "c"]
  empty_array    = []

  multiline_array_variable = [
    "d",
    "e",
    "f",
  ]

  unevaluateable_reference = var.name
  unevaluateable_interpolation = "this-${local.reference}"

  map_variable = {
    bool_variable   = true
    int_variable    = 1
    string_variable = "string"

    # Comment
    array_variable = ["a", "b", "c"]

    multiline_array_variable = [
      "d",
      "e",
      "f",
    ]
  }
}
