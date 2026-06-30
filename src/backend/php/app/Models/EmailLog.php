<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

/**
 * Eloquent model for the email_logs audit table.
 *
 * @property int              $id
 * @property string           $recipient
 * @property string           $subject
 * @property string           $body
 * @property string           $status       pending | sent | failed
 * @property \Carbon\Carbon   $created_at
 * @property \Carbon\Carbon|null $sent_at
 * @property string|null      $error_message
 */
class EmailLog extends Model
{
    /**
     * No Eloquent-managed updated_at column on this table.
     * created_at is handled by the DB default (useCurrent).
     */
    public $timestamps = false;

    protected $table = 'email_logs';

    protected $fillable = [
        'recipient',
        'subject',
        'body',
        'status',
        'sent_at',
        'error_message',
    ];

    protected $casts = [
        'created_at' => 'datetime',
        'sent_at'    => 'datetime',
    ];
}
