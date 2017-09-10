import 'package:angular/angular.dart';

@Pipe('mailto')
class MailtoPipe extends PipeTransform {
  String transform(String email) {
    return 'mailto:${email}';
  }
}
