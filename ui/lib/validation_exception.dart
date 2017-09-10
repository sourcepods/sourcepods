class ValidationException implements Exception {
  final message;

  ValidationException(this.message);

  String toString() {
    return "Validation Exception: $message";
  }
}
