// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {ThemeDefinition} from 'vuetify';
import '../styles/themes/default/main.scss';
import {black, blue, green, grey, neutral, orange, red, white, yellow} from './Colors';

const dark: ThemeDefinition = {
  dark: true,
  colors: {
    mbti: yellow[50],
    primary: blue[50],
    secondary: white,
    error: red[50],
    success: green[50],
    warning: yellow[50],
    info: blue[50],
    button: neutral[85],
    font: white,
    sidebar: neutral[30],

    // project status
    projectNew: orange[42],
    projectActive: blue[88],
    projectArchived: grey[30],
    projectGroup: orange[50],
    projectDeprecated: grey[30],

    // policy rule status
    prActive: blue[50],
    prInactive: grey[30],
    prDeprecated: orange[50],

    // version status
    versionNew: green[50],
    versionUnreviewed: grey[85],
    versionRejected: red[55],
    versionFreezed: grey[30],
    versionAcceptableAfterChanges: orange[65],
    versionApproved: green[50],

    borderCard: grey[30],

    tableRowOddBackgroundColor: neutral[20],
    tableRowOddBackgroundColorHover: neutral[35],
    tableRowEvenBackgroundColor: grey[5],
    tableRowEvenBackgroundColorHover: neutral[25],
    tableHeaderBackgroundColor: grey[5],
    tableBorderColor: grey[20],

    // modals
    successModalPlaceHolderBorderColor: green[70],
    errorModalPlaceHolderBorderColor: red[40],
    warningModalPlaceHolderBorderColor: yellow[70],
    warningPulseColor: yellow[70],
    modalBorderColor: neutral[60],

    // dialog
    dialogBackground: grey[5],
    dialogBorder: neutral[65],

    // markdown
    markdownBackground: grey[10],

    // licences
    licenceChartIcon: yellow[40],
    licenceForbidden: red[50],
    licenceNotApproved: neutral[75],

    //RuleButton
    ruleButton: neutral[80],

    //labels
    labels: neutral[80],
    // cards
    cardBackground: grey[10],
    cardBorder: grey[20],
    // background
    backgroundColor: grey[5],

    // PolicyStatus/Components Status
    policyStatusDeniedColor: red[50],
    policyStatusUnassertedColor: red[50],
    policyStatusWarnedColor: yellow[50],

    errorMessageColor: red[55],

    // Charts
    chartFLRed: blue[30],
    chartFLOrange: blue[35],
    chartFLYellow: blue[75],
    chartFLGreen: green[60],
    chartFLGrey: grey[50],

    chartRed: red[50],
    chartOrange: orange[50],
    chartYellow: yellow[50],
    chartGreen: green[50],
    chartGrey: neutral[65],
    chartLabelColor: white,

    tabsBackground: grey[10],
    tabsHeader: neutral[40],
    inactiveTabColor: grey[65],
    backgroundActiveTab: neutral[45],

    inputColor: grey[70],

    // Background Color of AND and OR
    backgroundColorAndOr: neutral[45],

    dashboardcardBackground: neutral[40],

    // Notification Bar
    notificationBarBackground: yellow[40],
    notificationBarTextColor: black,

    // fontColorTableSecond
    fontColorTableSecond: neutral[60],

    // approval status colors
    approvalApproved: green[50],
    approvalDeclined: red[50],
    approvalPending: orange[50],
  },
};
export default dark;
