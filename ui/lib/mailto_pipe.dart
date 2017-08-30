import 'package:angular2/angular2.dart';

@Pipe('mailto')
class MailtoPipe extends PipeTransform {
  String transform(String email) {
    return 'mailto:${email}';
  }
}
