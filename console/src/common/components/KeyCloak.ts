import {KcConfig} from '@/common/config/EnvConfig';
import Keycloak from 'keycloak-js';

// Setup Keycloak instance as needed
// Pass initialization options as required or leave blank to load from 'keycloak.json'
const keycloak = new Keycloak(KcConfig);

export default keycloak;
