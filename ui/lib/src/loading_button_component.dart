import 'package:angular/angular.dart';

@Component(
  selector: 'gitpods-loading-button',
  directives: [coreDirectives],
  template: '''
    <button class="uk-button uk-button-primary" type="submit" [disabled]="loading">
      <span [class.hide]="loading">{{ text }}</span>
      <div uk-spinner *ngIf="loading"></div>
    </button>
  ''',
  styles: [
    """
    button {
      position: relative;
    }
    .uk-spinner {
      position: absolute;
      top: 4px;
      left: 0;
      right: 0;
    }
    span.hide {
    opacity: 0;
    }
  """
  ],
)
class LoadingButtonComponent {
  @Input()
  String text;

  @Input()
  bool loading;
}
