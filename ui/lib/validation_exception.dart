class ValidationException implements Exception {
  final message;

  ValidationException(this.message);

  @override
  String toString() {
    return "Validation Exception: $message";
  }
}
