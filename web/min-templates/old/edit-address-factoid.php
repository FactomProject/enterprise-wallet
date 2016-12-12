<?php

$pageTitle = 'Edit Factoid Address';
$mainClass = 'send-factoids';
$activeNav = 2;

ob_start(); ?>

<div class="row">
                        <div class="columns">
                            <a href="/address-book.php" class="button close-button" data-close aria-label="Close reveal"><span aria-hidden="true">&times;</span></a>
                            <h1>Edit Factoid Address</h1>
                            <form>
                                <div class="row">
                                    <div class="medium-8 columns">
                                        <label>Public key:</label>
                                    </div>
                                    <div class="medium-4 columns">
                                        <label>Nickname:</label>
                                    </div>
                                </div>
                                <div class="row">
                                    <div class="medium-8 columns">
                                        <input type="text" disabled="true" name="public-key" value="FA2iH4u6cvTEac94PggP66rtM61NS3mURshXcPH6UwUQ5Mpoxch8">
                                    </div>
                                    <div class="medium-4 columns">
                                        <input type="text" name="alias" value="factoid1">
                                    </div>
                                </div>
                                <div class="row">
                                    <div class="medium-8 hide-for-small-only columns">

                                    </div>
                                    <div class="medium-4 small-12 columns">
                                        <a href="#" class="button expanded">Save Changes</a>
                                    </div>
                                </div>
                                <hr>
                                <div class="row collapse">
                                    <div class="large-4 medium-6 small-12 columns">
                                        <a href="#" class="button expanded">Display Private Key</a>
                                    </div>
                                    <div class="large-8 medium-6 small-12 columns">
                                        <input type="text" name="private-key" disabled="true" value="">
                                    </div>
                                </div>
                            </form>
                        </div>
                    </div>

<?php

$mainContent = ob_get_contents();
ob_end_clean();

include 'template.php';

?>