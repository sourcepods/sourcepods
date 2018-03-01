import 'package:angular/angular.dart';
import 'package:crypto/crypto.dart';

@Component(
    selector: 'gravatar',
    template: '<img class="{{classs}}" src="https://www.gravatar.com/avatar/{{hash}}?s={{size}}&d=mm&r=g"/>'
)
class GravatarComponent implements OnInit {
  @Input()
  String email;

  @Input()
  String size = '64';

  @Input('class')
  String classs = '';

  Digest hash;

  @override
  void ngOnInit() {
    this.hash = md5.convert(this.email.codeUnits);
  }
}
