import {Rights} from '@disclosure-portal/model/Rights';
import {UserDto} from '@disclosure-portal/model/Users';

export default interface SimpleProfileData {
  rights: Rights;
  profile: UserDto;
  allowed: boolean;
}
