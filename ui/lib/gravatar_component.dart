import 'package:angular/angular.dart';
import 'package:crypto/crypto.dart';


@Component(
    selector: 'gravatar',
    template: '''Bank Name: {{email}} - Account Id: {{size}}''')
class Gravatar {
  @Input()
  String email;

  @Input()
  String size;

//  Gravatar(String this.email, String this.size);
}

//@Component(
//  selector: 'gravatar',
//  template: '<img src="https://www.gravatar.com/avatar/{{hash}}?s={{size}}&d=mm&r=g" alt="">',
//)
//class Gravatar implements OnInit {
//  Digest hash;
//
//  @Attribute('email')
//  String email;
//  @Attribute('size')
//  int size = 128;
//
//  Gravatar(this.email, this.size);
//
//  @override
//  ngOnInit() {
//    print('OnInit: ${this.email}');
//    print('${this.email} - ${size}');
//    hash = md5.convert('mail@matthiasloibl.com'.codeUnits);
//  }
//}

//@Directive(selector: 'input')
//class InputAttrDirective {
//  InputAttrDirective(@Attribute('type') String type) {
//    // type would be 'text' in this example
//  }
//}

