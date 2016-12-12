<?php

$pageTitle = 'Receive Factoids';
$mainClass = 'receive-factoids';
$activeNav = 1;

ob_start(); ?>

<div class="row">
                        <div class="columns">
                            <a href="/" class="button close-button" data-close aria-label="Close reveal"><span aria-hidden="true">&times;</span></a>
                            <h1>Receive Factoids</h1>
                            <div class="row">
                                <div class="small-12 large-4 columns">
                                    <img src="img/qr_receiving.png" class="qr-receiving"><br><br>
                                </div>
                                <div class="small-12 large-8 columns">
                                    <form>
                                        <label>Receiving address: <small><a href="#">Manage Addresses</a></small></label>
                                        <select id="receiving-address" name="receiving-address">
                                            <option value="factoid1">factoid1 (FA22iH4u8cvTEac94PggP66rtM61NS3mURshXcPH6UwUQ5Mpoxch8)</option>
                                            <option value="factoid2">factoid2 (FA22iH4u8cvTEac94PggP66rtM61NS3mURshXcPH6UwUQ5Mpoxch8)</option>
                                            <option value="factoid3">factoid3 (FA22iH4u8cvTEac94PggP66rtM61NS3mURshXcPH6UwUQ5Mpoxch8)</option>
                                            <option value="factoid4">factoid5 (FA22iH4u8cvTEac94PggP66rtM61NS3mURshXcPH6UwUQ5Mpoxch8)</option>
                                            <option value="factoid5">factoid6 (FA22iH4u8cvTEac94PggP66rtM61NS3mURshXcPH6UwUQ5Mpoxch8)</option>
                                            <option value="factoid6">factoid7 (FA22iH4u8cvTEac94PggP66rtM61NS3mURshXcPH6UwUQ5Mpoxch8)</option>
                                        </select>
                                    </form>
                                    <div class="row">
                                        <div class="columns text-right">
                                            <a href="#" class="button copy-address">Copy to Clipboard</a>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>

<?php

$mainContent = ob_get_contents();
ob_end_clean();

include 'template.php';

?>